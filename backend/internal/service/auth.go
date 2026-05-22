package service

import (
	"backend/internal/config"
	"backend/internal/crypto"
	"backend/internal/logging"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"crypto/subtle"
	"errors"
	"time"
)

var (
	ErrInvalidAdminCredentials  = errors.New("invalid admin credentials")
	ErrInvalidClientCredentials = errors.New("invalid device credentials")
	ErrInvalidRegistrationCode  = errors.New("invalid registration code")
)

const (
	RegistrationExpiryDuration = 10 * time.Minute
	TokenExpiry                = 40 * 24 * time.Hour
)

type AuthService struct {
	config config.AuthConfig
	repo   repository.Repository
}

func NewAuthService(cfg config.AuthConfig, repo repository.Repository) *AuthService {
	return &AuthService{
		config: cfg,
		repo:   repo,
	}
}

func (s *AuthService) ValidateAdminToken(apiKey string) error {
	expectedKey := s.config.AdminAPIKey
	if subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedKey)) != 1 {
		return ErrInvalidAdminCredentials
	}

	return nil
}

func (s *AuthService) ValidateDeviceToken(ctx context.Context, token string) (models.Device, error) {
	tokenHash := crypto.HashToken(token)

	device, err := s.repo.GetDeviceByToken(ctx, tokenHash)
	if err != nil {
		return models.Device{}, err
	}

	if device == nil || device.IsExpired() {
		return models.Device{}, ErrInvalidClientCredentials
	}

	// Update last seen timestamp
	if err := s.repo.UpdateDeviceLastSeenByToken(ctx, tokenHash, time.Now()); err != nil {
		return models.Device{}, err
	}

	return *device, nil
}

func (s *AuthService) RegisterDevice(ctx context.Context, registrationCode string, deviceInfo string) (string, error) {
	logger := logging.FromContext(ctx)

	code, err := s.repo.GetRegistrationCode(ctx)
	if err != nil {
		return "", err
	}

	if code == nil || code.IsExpired() || subtle.ConstantTimeCompare([]byte(registrationCode), []byte(code.Code)) != 1 {
		logger.DebugContext(ctx, "invalid registration code attempt", "provided_code", registrationCode)
		return "", ErrInvalidRegistrationCode
	}

	secureToken := crypto.GenerateSecureToken(32)
	tokenHash := crypto.HashToken(secureToken)

	expiry := time.Now().Add(TokenExpiry)
	err = s.repo.WithTx(ctx, nil, func(repo repository.Repository) error {
		if err := repo.DeleteRegistrationCode(ctx); err != nil {
			return err
		}

		if err := repo.RegisterDevice(ctx, tokenHash, deviceInfo, expiry); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.ErrorContext(ctx, "failed to register device", "error", err)
		return "", err
	}

	logger.InfoContext(ctx, "successfully registered device", "device_info", deviceInfo)
	return secureToken, nil
}

func (s *AuthService) CreateRegistrationCode(ctx context.Context) (models.RegistrationCode, error) {
	createdAt := time.Now()
	expiry := createdAt.Add(RegistrationExpiryDuration)

	code, err := crypto.GenerateHumanReadableCode()
	if err != nil {
		return models.RegistrationCode{}, err
	}

	model := models.RegistrationCode{
		Code:      code,
		ExpiresAt: expiry,
		CreatedAt: createdAt,
	}

	if err := s.repo.RefreshRegistrationCode(ctx, model); err != nil {
		return model, err
	}

	return model, nil
}

func (s *AuthService) InvalidateRegistrationCode(ctx context.Context) error {
	return s.repo.DeleteRegistrationCode(ctx)
}

func (s *AuthService) GetRegistrationCode(ctx context.Context) (*models.RegistrationCode, error) {
	return s.repo.GetRegistrationCode(ctx)
}

func (s *AuthService) ListDevices(ctx context.Context) ([]models.Device, error) {
	return s.repo.ListDevices(ctx)
}

func (s *AuthService) RevokeDevice(ctx context.Context, deviceID string) error {
	return s.repo.DeleteDevice(ctx, deviceID)
}

func (s *AuthService) UpdateDeviceInfo(ctx context.Context, deviceID string, deviceInfo string) error {
	return s.repo.UpdateDeviceInfo(ctx, deviceID, deviceInfo)
}
