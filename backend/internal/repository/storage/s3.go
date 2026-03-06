package storage

import (
	"backend/internal/config"
	"backend/internal/crypto"
	"context"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type listObject struct {
	Key string `xml:"Key"`
}

type listBucketResult struct {
	XMLName               xml.Name     `xml:"ListBucketResult"`
	Contents              []listObject `xml:"Contents"`
	IsTruncated           bool         `xml:"IsTruncated"`
	NextContinuationToken string       `xml:"NextContinuationToken"`
}

type deleteObject struct {
	XMLName xml.Name `xml:"Object"`
	Key     string   `xml:"Key"`
}

type deleteRequest struct {
	XMLName xml.Name       `xml:"Delete"`
	Quiet   bool           `xml:"Quiet"`
	Objects []deleteObject `xml:"Object"`
}

type deleteErrorEntry struct {
	Key     string `xml:"Key"`
	Code    string `xml:"Code"`
	Message string `xml:"Message"`
}

type deleteResult struct {
	XMLName xml.Name           `xml:"DeleteResult"`
	Errors  []deleteErrorEntry `xml:"Error"`
}

const maxDeleteBatch = 1000

type S3FileStore struct {
	endpoint string
	bucket   string
	signer   *crypto.SigV4Signer
	client   *http.Client
}

func NewS3FileStore(conf config.S3FileStorageConfig) *S3FileStore {
	signer := crypto.NewSigV4Signer(conf.Region, conf.AccessKey, conf.SecretKey)
	return &S3FileStore{
		endpoint: conf.Endpoint,
		bucket:   conf.Bucket,
		signer:   signer,
		client: &http.Client{
			Timeout: 30 * 1e9, // 30 seconds
		},
	}
}

func (s *S3FileStore) PresignGet(_ context.Context, key string, ttl time.Duration) (string, error) {
	req, err := http.NewRequest("GET", s.endpoint+"/"+s.bucket+"/"+key, nil)
	if err != nil {
		return "", err
	}

	s.signer.PresignRequest(req, ttl)
	return req.URL.String(), nil
}

func (s *S3FileStore) Put(ctx context.Context, filename string, r io.ReadSeeker) (string, error) {
	key := s.generateKey(filename)

	hasher := sha256.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return "", err
	}
	payloadHash := fmt.Sprintf("%x", hasher.Sum(nil))
	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", s.endpoint+"/"+s.bucket+"/"+key, r)
	if err != nil {
		return "", err
	}

	s.signer.SignRequest(req, payloadHash)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	return key, nil
}

func (s *S3FileStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.endpoint+"/"+s.bucket+"/"+key, nil)
	if err != nil {
		return nil, err
	}

	s.signer.SignRequest(req, "")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to get file: %s", resp.Status)
	}

	return resp.Body, nil
}

func (s *S3FileStore) Delete(ctx context.Context, key string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", s.endpoint+"/"+s.bucket+"/"+key, nil)
	if err != nil {
		return err
	}

	s.signer.SignRequest(req, "")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete file: %s", resp.Status)
	}

	return nil
}

func (s *S3FileStore) DeleteMany(ctx context.Context, keys []string) error {
	for i := 0; i < len(keys); i += maxDeleteBatch {
		end := min(i+maxDeleteBatch, len(keys))
		if err := s.deleteBatch(ctx, keys[i:end]); err != nil {
			return err
		}
	}
	return nil
}

func (s *S3FileStore) List(ctx context.Context) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		var token string
		for {
			u := s.endpoint + "/" + s.bucket + "?list-type=2"
			if token != "" {
				u += "&continuation-token=" + url.QueryEscape(token)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
			if err != nil {
				yield("", err)
				return
			}

			s.signer.SignRequest(req, "")

			resp, err := s.client.Do(req)
			if err != nil {
				yield("", err)
				return
			}

			var result listBucketResult
			err = xml.NewDecoder(resp.Body).Decode(&result)
			resp.Body.Close()
			if err != nil {
				yield("", fmt.Errorf("s3 list: decode response: %w", err))
				return
			}

			for _, obj := range result.Contents {
				if !yield(obj.Key, nil) {
					return
				}
			}

			if !result.IsTruncated {
				return
			}
			token = result.NextContinuationToken
		}
	}
}

func (s *S3FileStore) deleteBatch(ctx context.Context, keys []string) error {
	payload := deleteRequest{Quiet: true}
	for _, k := range keys {
		payload.Objects = append(payload.Objects, deleteObject{Key: k})
	}

	body, err := xml.Marshal(payload)
	if err != nil {
		return fmt.Errorf("s3 delete many: marshal: %w", err)
	}

	bodyHash := fmt.Sprintf("%x", sha256.Sum256(body))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint+"/"+s.bucket+"?delete", strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("s3 delete many: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/xml")

	s.signer.SignRequest(req, bodyHash)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("s3 delete many: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("s3 delete many: status %d: %s", resp.StatusCode, respBody)
	}

	var result deleteResult
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Quiet mode with 200 OK — body may be empty, which is fine.
		return nil
	}

	if len(result.Errors) > 0 {
		e := result.Errors[0]
		return fmt.Errorf("s3 delete many: %s: %s (%s)", e.Key, e.Message, e.Code)
	}

	return nil
}

func (s *S3FileStore) generateKey(filename string) string {
	ext := filepath.Ext(filename)
	return crypto.UUIDV4() + ext
}
