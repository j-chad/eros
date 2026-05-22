package webpush

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

func GenerateVapidKeys() (private string, public string, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	privateBytes, err := privateKey.Bytes()
	if err != nil {
		return "", "", fmt.Errorf("failed to serialize private key: %w", err)
	}

	publicBytes, err := privateKey.PublicKey.Bytes()
	if err != nil {
		return "", "", fmt.Errorf("failed to serialize public key: %w", err)
	}

	private = base64.RawURLEncoding.EncodeToString(privateBytes)
	public = base64.RawURLEncoding.EncodeToString(publicBytes)
	return private, public, nil
}

func vapidAuthorization(endpoint string, privateKey *ecdsa.PrivateKey, subject string) (string, error) {
	jwt, err := buildVapidJWT(privateKey, endpoint, subject)
	if err != nil {
		return "", fmt.Errorf("build VAPID JWT: %w", err)
	}

	publicKeyBytes, err := privateKey.PublicKey.Bytes()
	if err != nil {
		return "", fmt.Errorf("serialize public key: %w", err)
	}
	publicKey := base64.RawURLEncoding.EncodeToString(publicKeyBytes)

	authorization := bytes.NewBufferString("vapid t=")
	authorization.WriteString(jwt)
	authorization.WriteString(",k=")
	authorization.WriteString(publicKey)

	return authorization.String(), nil
}

func buildVapidJWT(privateKey *ecdsa.PrivateKey, endpoint string, subject string) (string, error) {
	header, err := b64JSON(map[string]any{
		"typ": "JWT",
		"alg": "ES256",
	})
	if err != nil {
		return "", fmt.Errorf("build JWT header: %w", err)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse endpoint URL: %w", err)
	}
	aud := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	body, err := b64JSON(map[string]any{
		"aud": aud,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"sub": subject,
	})
	if err != nil {
		return "", fmt.Errorf("build JWT payload: %w", err)
	}

	signingInput := []byte(header + "." + body)
	hash := sha256.Sum256(signingInput)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("sign JWT payload: %w", err)
	}

	// zero-pad r & s to 32 bytes each
	signature := make([]byte, 64)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	copy(signature[32-len(rBytes):32], rBytes)
	copy(signature[64-len(sBytes):64], sBytes)

	return fmt.Sprintf("%s.%s.%s", header, body, base64.RawURLEncoding.EncodeToString(signature)), nil
}

func b64JSON(m map[string]any) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
