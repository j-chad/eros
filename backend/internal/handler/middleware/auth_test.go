package middleware

import (
	"backend/internal/testutil"
	"testing"
)

func TestParseAuthHeader_Admin(t *testing.T) {
	authType, token, ok := parseAuthHeader("Admin my-secret-key")
	testutil.True(t, ok, "should parse Admin header")
	testutil.Equal(t, authType, "Admin")
	testutil.Equal(t, token, "my-secret-key")
}

func TestParseAuthHeader_Bearer(t *testing.T) {
	authType, token, ok := parseAuthHeader("Bearer abc123")
	testutil.True(t, ok, "should parse Bearer header")
	testutil.Equal(t, authType, "Bearer")
	testutil.Equal(t, token, "abc123")
}

func TestParseAuthHeader_Empty(t *testing.T) {
	_, _, ok := parseAuthHeader("")
	testutil.False(t, ok, "empty header should fail")
}

func TestParseAuthHeader_MissingToken(t *testing.T) {
	_, _, ok := parseAuthHeader("Bearer")
	testutil.False(t, ok, "header without token should fail")
}

func TestParseAuthHeader_UnknownScheme(t *testing.T) {
	_, _, ok := parseAuthHeader("Basic abc123")
	testutil.False(t, ok, "unsupported scheme should fail")
}

func TestParseAuthHeader_ExtraWhitespace(t *testing.T) {
	_, _, ok := parseAuthHeader("Bearer token extra")
	testutil.True(t, ok, "Sscanf reads first two tokens, extra ignored")
}

func TestParseAuthHeader_OnlyScheme(t *testing.T) {
	_, _, ok := parseAuthHeader("Admin")
	testutil.False(t, ok, "scheme without token should fail")
}
