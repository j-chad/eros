package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
