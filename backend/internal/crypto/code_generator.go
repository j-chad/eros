package crypto

import (
	"crypto/rand"
	"encoding/base64"
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

func GenerateSecureToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Encode as URL-safe base64 without padding
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
