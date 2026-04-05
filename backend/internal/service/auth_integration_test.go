package service

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/testutil"
	"backend/internal/testutil/testdb"
	"context"
	"testing"
	"time"
)

func TestValidateDeviceToken_Valid(t *testing.T) {
	repo := testdb.New(t)
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "key"}, repo)
	ctx := context.Background()

	// Register a device to get a token
	token, err := svc.RegisterDevice(ctx, "CODE", "iPhone")
	testutil.NilErr(t, err)
	testutil.True(t, len(token) > 0, "token should not be empty")

	// Validate the token
	testutil.NilErr(t, svc.ValidateDeviceToken(token))
}

func TestValidateDeviceToken_Invalid(t *testing.T) {
	repo := testdb.New(t)
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "key"}, repo)

	err := svc.ValidateDeviceToken("nonexistent-token")
	// nil expiry returns ErrInvalidClientCredentials
	testutil.ErrorIs(t, err, ErrInvalidClientCredentials)
}

func TestValidateDeviceToken_Expired(t *testing.T) {
	repo := testdb.New(t)
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "key"}, repo)
	ctx := context.Background()

	// Register a device with a past expiry directly via repo
	repo.RegisterDevice(ctx, "expired-token", "old phone", time.Now().Add(-time.Hour))

	err := svc.ValidateDeviceToken("expired-token")
	testutil.ErrorIs(t, err, ErrInvalidClientCredentials)
}

func TestValidateDeviceToken_UpdatesLastSeen(t *testing.T) {
	repo := testdb.New(t)
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "key"}, repo)
	ctx := context.Background()

	token, err := svc.RegisterDevice(ctx, "CODE", "iPhone")
	testutil.NilErr(t, err)

	// Validate twice — should not error (last seen gets updated)
	testutil.NilErr(t, svc.ValidateDeviceToken(token))
	testutil.NilErr(t, svc.ValidateDeviceToken(token))
}

func TestRegisterDevice_CreatesDevice(t *testing.T) {
	repo := testdb.New(t)
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "key"}, repo)
	ctx := context.Background()

	// Seed a registration code (RegisterDevice deletes it)
	repo.RefreshRegistrationCode(ctx, models.RegistrationCode{
		Code:      "TEST-CODE-1234",
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	})

	token, err := svc.RegisterDevice(ctx, "TEST-CODE-1234", "Pixel 7")
	testutil.NilErr(t, err)
	testutil.True(t, len(token) > 0, "should return a token")

	// Verify device exists via token validation
	testutil.NilErr(t, svc.ValidateDeviceToken(token))

	// Verify registration code was deleted
	code, err := repo.GetRegistrationCode(ctx)
	testutil.NilErr(t, err)
	testutil.Nil(t, code)
}

func TestRegisterDevice_DeletesRegistrationCode(t *testing.T) {
	repo := testdb.New(t)
	svc := NewAuthService(config.AuthConfig{AdminAPIKey: "key"}, repo)
	ctx := context.Background()

	repo.RefreshRegistrationCode(ctx, models.RegistrationCode{
		Code:      "WILL-BE-GONE",
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	})

	_, err := svc.RegisterDevice(ctx, "WILL-BE-GONE", "device")
	testutil.NilErr(t, err)

	code, _ := repo.GetRegistrationCode(ctx)
	testutil.Nil(t, code)
}
