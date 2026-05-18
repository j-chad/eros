package crypto

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	aws4Request     = "aws4_request"
	signingAlgo     = "AWS4-HMAC-SHA256"
	unsignedPayload = "UNSIGNED-PAYLOAD"
	timeFormat      = "20060102T150405Z"
	dateFormat      = "20060102"
	service         = "s3"
)

// SigV4Signer implements AWS Signature Version 4 signing for S3 requests.
// It can sign both regular HTTP requests and generate pre-signed URLs.
type SigV4Signer struct {
	region    string
	accessKey string
	secretKey string
}

func NewSigV4Signer(region, accessKey, secretKey string) *SigV4Signer {
	return &SigV4Signer{
		region:    region,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

func (s *SigV4Signer) PresignRequest(req *http.Request, expires time.Duration) {
	t := time.Now().UTC()
	credentialScope := fmt.Sprintf("%s/%s/%s/%s", t.Format(dateFormat), s.region, service, aws4Request)

	// For presigned URLs the auth parameters go in the query string and must be
	// part of the canonical request *before* signing.
	query := req.URL.Query()
	query.Set("X-Amz-Algorithm", signingAlgo)
	query.Set("X-Amz-Credential", fmt.Sprintf("%s/%s", s.accessKey, credentialScope))
	query.Set("X-Amz-Date", t.Format(timeFormat))
	query.Set("X-Amz-Expires", fmt.Sprintf("%d", int(expires.Seconds())))
	query.Set("X-Amz-SignedHeaders", "host")
	req.URL.RawQuery = query.Encode()

	// Ensure Host header is present for signing.
	if req.Header.Get("Host") == "" {
		req.Header.Set("Host", req.URL.Host)
	}

	signedHeaders, canonicalHeaders := s.buildCanonicalHeaders(req)
	canonicalReq := s.buildCanonicalRequest(req, unsignedPayload, signedHeaders, canonicalHeaders)
	stringToSign := s.buildStringToSign(canonicalReq, credentialScope, t)
	signingKey := s.deriveSigningKey(t)
	signature := fmt.Sprintf("%x", hmacSHA256(signingKey, stringToSign))

	query.Set("X-Amz-Signature", signature)
	req.URL.RawQuery = query.Encode()
}

func (s *SigV4Signer) SignRequest(req *http.Request, payloadHash string) {
	if payloadHash == "" {
		payloadHash = unsignedPayload
	}

	t := time.Now().UTC()

	req.Header.Set("x-amz-date", t.Format(timeFormat))
	req.Header.Set("x-amz-content-sha256", payloadHash)
	if req.Header.Get("Host") == "" {
		req.Header.Set("Host", req.Host)
	}

	credentialScope := fmt.Sprintf("%s/%s/%s/%s", t.Format(dateFormat), s.region, service, aws4Request)

	signedHeaders, canonicalHeaders := s.buildCanonicalHeaders(req)
	canonicalReq := s.buildCanonicalRequest(req, payloadHash, signedHeaders, canonicalHeaders)
	stringToSign := s.buildStringToSign(canonicalReq, credentialScope, t)
	signingKey := s.deriveSigningKey(t)
	signature := fmt.Sprintf("%x", hmacSHA256(signingKey, stringToSign))

	authHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		signingAlgo, s.accessKey, credentialScope, signedHeaders, signature)
	req.Header.Set("Authorization", authHeader)
}

func (s *SigV4Signer) deriveSigningKey(t time.Time) []byte {
	dateKey := hmacSHA256([]byte("AWS4"+s.secretKey), []byte(t.Format(dateFormat)))
	regionKey := hmacSHA256(dateKey, []byte(s.region))
	serviceKey := hmacSHA256(regionKey, []byte(service))
	signingKey := hmacSHA256(serviceKey, []byte(aws4Request))
	return signingKey
}

func (s *SigV4Signer) buildStringToSign(canonicalReq, scope string, t time.Time) []byte {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(canonicalReq)))
	return []byte(signingAlgo + "\n" +
		t.Format(timeFormat) + "\n" +
		scope + "\n" +
		hash)
}

func (s *SigV4Signer) buildCanonicalRequest(req *http.Request, payloadHash, signedHeaders, canonicalHeaders string) string {
	if req.Method == "" {
		req.Method = "GET"
	}

	return req.Method + "\n" +
		s.awsURIEncode(req.URL.Path, false) + "\n" +
		s.encodeQuery(req.URL.Query()) + "\n" +
		canonicalHeaders + "\n\n" +
		signedHeaders + "\n" +
		payloadHash
}

func (s *SigV4Signer) buildCanonicalHeaders(req *http.Request) (signed string, canonical string) {
	type hdr struct {
		key, val string
	}

	var headers []hdr
	for k, v := range req.Header {
		lk := strings.ToLower(k)
		if lk == "host" || strings.HasPrefix(lk, "x-amz-") || lk == "content-type" || lk == "content-length" {
			headers = append(headers, hdr{lk, strings.TrimSpace(v[0])})
		}
	}
	sort.Slice(headers, func(i, j int) bool { return headers[i].key < headers[j].key })

	var signedParts, canonicalParts []string
	for _, h := range headers {
		signedParts = append(signedParts, h.key)
		canonicalParts = append(canonicalParts, h.key+":"+h.val)
	}

	return strings.Join(signedParts, ";"), strings.Join(canonicalParts, "\n")
}

// encodeQuery encodes URL query parameters in a canonical form for signing.
func (s *SigV4Signer) encodeQuery(vals url.Values) string {
	if len(vals) == 0 {
		return ""
	}
	keys := make([]string, 0, len(vals))
	for k := range vals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		ek := s.awsURIEncode(k, true)
		for _, v := range vals[k] {
			parts = append(parts, ek+"="+s.awsURIEncode(v, true))
		}
	}
	return strings.Join(parts, "&")
}

// awsURIEncode encodes a string for use in a URI, following AWS's specific rules.
func (s *SigV4Signer) awsURIEncode(path string, encodeSlash bool) string {
	var b strings.Builder
	for _, c := range []byte(path) {
		switch {
		case (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '_' || c == '-' || c == '~' || c == '.':
			b.WriteByte(c)
		case c == '/' && !encodeSlash:
			b.WriteByte('/')
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}
