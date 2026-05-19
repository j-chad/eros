package crypto

import (
	"backend/internal/testutil"
	"regexp"
	"strings"
	"testing"
)

func TestUUIDV4_Format(t *testing.T) {
	id := UUIDV4()
	testutil.True(t, UUIDV4Regex.MatchString(id), "should match UUID v4 regex: "+id)
}

func TestUUIDV4_VersionBit(t *testing.T) {
	id := UUIDV4()
	testutil.Equal(t, id[14], byte('4'))
}

func TestUUIDV4_VariantBit(t *testing.T) {
	id := UUIDV4()
	variant := id[19]
	testutil.True(t, variant == '8' || variant == '9' || variant == 'a' || variant == 'b',
		"variant nibble should be one of [89ab], got "+string(variant))
}

func TestUUIDV4_Unique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := UUIDV4()
		testutil.False(t, seen[id], "duplicate UUID: "+id)
		seen[id] = true
	}
}

func TestUUIDV4_Length(t *testing.T) {
	testutil.Equal(t, len(UUIDV4()), 36)
}

func TestGenerateHumanReadableCode_Format(t *testing.T) {
	code, err := GenerateHumanReadableCode()
	testutil.NilErr(t, err)

	pattern := regexp.MustCompile(`^[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}$`)
	testutil.True(t, pattern.MatchString(code), "code should match XXXX-XXXX-XXXX: "+code)
}

func TestGenerateHumanReadableCode_Unique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code, err := GenerateHumanReadableCode()
		testutil.NilErr(t, err)
		testutil.False(t, seen[code], "duplicate code: "+code)
		seen[code] = true
	}
}

func TestGenerateSecureToken_Length(t *testing.T) {
	token := GenerateSecureToken(32)
	// 32 bytes base64url-encoded without padding = 43 chars
	testutil.Equal(t, len(token), 43)
}

func TestGenerateSecureToken_NoPadding(t *testing.T) {
	token := GenerateSecureToken(32)
	testutil.False(t, strings.Contains(token, "="), "should not contain padding")
}

func TestGenerateSecureToken_URLSafe(t *testing.T) {
	for i := 0; i < 100; i++ {
		token := GenerateSecureToken(32)
		testutil.False(t, strings.ContainsAny(token, "+/"), "should be URL-safe: "+token)
	}
}

func TestHashToken_Deterministic(t *testing.T) {
	token := "test-token-value"
	testutil.Equal(t, HashToken(token), HashToken(token))
}

func TestHashToken_HexEncoded(t *testing.T) {
	hash := HashToken("test")
	// SHA-256 produces 32 bytes = 64 hex chars
	testutil.Equal(t, len(hash), 64)
	pattern := regexp.MustCompile(`^[0-9a-f]{64}$`)
	testutil.True(t, pattern.MatchString(hash), "should be lowercase hex: "+hash)
}

func TestHashToken_DifferentInputs(t *testing.T) {
	testutil.True(t, HashToken("a") != HashToken("b"), "different inputs should produce different hashes")
}
