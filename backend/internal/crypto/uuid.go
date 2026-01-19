package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// UUIDV4 generates a random UUID version 4 as per RFC 4122.
func UUIDV4() string {
	var uuid [16]byte
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		panic(err)
	}

	// Set version (4) and variant (RFC 4122)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // variant 10xxxxxx

	// Format: 8-4-4-4-12 hex chars
	var out [36]byte
	hex.Encode(out[0:8], uuid[0:4])
	out[8] = '-'
	hex.Encode(out[9:13], uuid[4:6])
	out[13] = '-'
	hex.Encode(out[14:18], uuid[6:8])
	out[18] = '-'
	hex.Encode(out[19:23], uuid[8:10])
	out[23] = '-'
	hex.Encode(out[24:36], uuid[10:16])

	return string(out[:])
}
