package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateHumanReadableCode generates a human-readable code string.
// In the format of XXXX-XXXX-XXXX, where X is an uppercase letter or digit.
// Total entropy of 62 bits.
func GenerateHumanReadableCode() (string, error) {
	result := make([]byte, 12)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}

	return fmt.Sprintf("%s-%s-%s", string(result[0:4]), string(result[4:8]), string(result[8:12])), nil
}

func RandomBytes(n int) []byte {
	result := make([]byte, n)
	_, _ = rand.Read(result)
	return result
}

func GenerateSecureToken(n int) string {
	bytes := RandomBytes(n)
	return base64.RawURLEncoding.EncodeToString(bytes)
}

// HashToken returns the SHA-256 hex digest of a token.
// Used to store and look up device tokens without keeping plaintext in the database.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
