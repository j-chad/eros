package service

import (
	"backend/internal/config"
	"backend/internal/crypto"
	"backend/internal/repository"
	"context"
	"crypto/subtle"
	"errors"
	"time"
)

var (
	ErrInvalidAdminCredentials = errors.New("invalid admin credentials")
	ErrInvalidRegistrationCode = errors.New("invalid registration code")
)

const TokenLength = 32
const TokenExpiry = 40 * 24 * time.Hour

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

func (s *AuthService) ValidateAdminAPIKey(apiKey string) error {
	expectedKey := s.config.AdminAPIKey
	if subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedKey)) != 1 {
		return ErrInvalidAdminCredentials
	}

	return nil
}

func (s *AuthService) RegisterDevice(ctx context.Context, registrationCode string, deviceInfo string) (string, error) {
	code, err := s.repo.GetRegistrationCode(ctx)
	if err != nil {
		return "", err
	}

	if code == nil || code.IsExpired() || subtle.ConstantTimeCompare([]byte(registrationCode), []byte(code.Code)) != 1 {
		return "", ErrInvalidRegistrationCode
	}

	secureToken, err := crypto.GenerateSecureToken(32)
	if err != nil {
		return "", err
	}

	expiry := time.Now().Add(TokenExpiry)
	err = s.repo.WithTx(ctx, func(repo repository.Repository) error {
		if err := repo.DeleteRegistrationCode(ctx); err != nil {
			return err
		}

		if err := repo.RegisterDevice(ctx, secureToken, deviceInfo, expiry); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return secureToken, nil
}
