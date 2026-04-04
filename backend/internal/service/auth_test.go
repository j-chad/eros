package service

import (
	"backend/internal/config"
	"backend/internal/testutil"
	"testing"
)

func TestValidateAdminToken_Correct(t *testing.T) {
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "my-secret"}, nil)
	testutil.NilErr(t, svc.ValidateAdminToken("my-secret"))
}

func TestValidateAdminToken_Wrong(t *testing.T) {
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "my-secret"}, nil)
	testutil.ErrorIs(t, svc.ValidateAdminToken("wrong"), ErrInvalidAdminCredentials)
}

func TestValidateAdminToken_Empty(t *testing.T) {
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "my-secret"}, nil)
	testutil.ErrorIs(t, svc.ValidateAdminToken(""), ErrInvalidAdminCredentials)
}

func TestValidateAdminToken_EmptyConfig(t *testing.T) {
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: ""}, nil)
	// Empty matches empty — constant time compare returns 1
	testutil.NilErr(t, svc.ValidateAdminToken(""))
}

func TestValidateAdminToken_TimingResistant(t *testing.T) {
	// Ensure it uses constant-time compare — different lengths should still error, not panic
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "short"}, nil)
	testutil.ErrorIs(t, svc.ValidateAdminToken("a-much-longer-key-that-differs"), ErrInvalidAdminCredentials)
}
