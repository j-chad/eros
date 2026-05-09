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
			Timeout: 30 * time.Second,
		},
	}
}

func (s *S3FileStore) PresignGet(_ context.Context, key string, ttl time.Duration) (string, error) {
	req, err := http.NewRequest("GET", s.objectURL(key), nil)
	if err != nil {
		return "", fmt.Errorf("s3 presign: %w", err)
	}

	s.signer.PresignRequest(req, ttl)
	return req.URL.String(), nil
}

func (s *S3FileStore) Put(ctx context.Context, filename, mime string, r io.ReadSeeker) (string, error) {
	key := s.generateKey(filename)

	hasher := sha256.New()
	size, err := io.Copy(hasher, r)
	if err != nil {
		return "", fmt.Errorf("s3 put: hash payload: %w", err)
	}
	payloadHash := fmt.Sprintf("%x", hasher.Sum(nil))
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("s3 put: reset reader: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.objectURL(key), r)
	if err != nil {
		return "", fmt.Errorf("s3 put: build request: %w", err)
	}
	req.ContentLength = size
	req.Header.Set("Content-Type", mime)

	s.signer.SignRequest(req, payloadHash)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("s3 put: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", s.readErrorResponse("put", resp)
	}

	return key, nil
}

func (s *S3FileStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.objectURL(key), nil)
	if err != nil {
		return nil, fmt.Errorf("s3 get: build request: %w", err)
	}

	s.signer.SignRequest(req, "")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("s3 get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, s.readErrorResponse("get", resp)
	}

	return resp.Body, nil
}

func (s *S3FileStore) Delete(ctx context.Context, key string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, s.objectURL(key), nil)
	if err != nil {
		return fmt.Errorf("s3 delete: build request: %w", err)
	}

	s.signer.SignRequest(req, "")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("s3 delete: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return s.readErrorResponse("delete", resp)
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
			u := s.bucketURL("list-type=2")
			if token != "" {
				u += "&continuation-token=" + url.QueryEscape(token)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
			if err != nil {
				yield("", fmt.Errorf("s3 list: build request: %w", err))
				return
			}

			s.signer.SignRequest(req, "")

			resp, err := s.client.Do(req)
			if err != nil {
				yield("", fmt.Errorf("s3 list: %w", err))
				return
			}

			var result listBucketResult
			err = xml.NewDecoder(resp.Body).Decode(&result)
			closeErr := resp.Body.Close()
			if closeErr != nil {
				yield("", fmt.Errorf("s3 list: read response: %w", err))
				return
			}
			if err != nil {
				yield("", s.readErrorResponse("list", resp))
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.bucketURL("delete"), strings.NewReader(string(body)))
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
		return s.readErrorResponse("delete many", resp)
	}

	var result deleteResult
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Quiet mode with 200 OK - body may be empty, which is fine.
		return nil
	}

	if len(result.Errors) > 0 {
		e := result.Errors[0]
		return fmt.Errorf("s3 delete many: %s: %s (%s)", e.Key, e.Message, e.Code)
	}

	return nil
}

func (s *S3FileStore) objectURL(key string) string {
	return s.endpoint + "/" + s.bucket + "/" + key
}

func (s *S3FileStore) bucketURL(query string) string {
	if query != "" {
		return s.endpoint + "/" + s.bucket + "?" + query
	}
	return s.endpoint + "/" + s.bucket
}

func (s *S3FileStore) generateKey(filename string) string {
	ext := filepath.Ext(filename)
	return crypto.UUIDV4() + ext
}

func (s *S3FileStore) readErrorResponse(op string, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("s3 %s: read error response: %w", op, err)
	}

	return fmt.Errorf("s3 %s: status %d: %s", op, resp.StatusCode, body)
}

// XML types for S3 API responses

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
