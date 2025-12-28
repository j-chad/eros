package service

import (
	"backend/internal/config"
	"backend/internal/repository"
	"crypto/subtle"
	"errors"
)

var (
	ErrInvalidAdminCredentials = errors.New("invalid admin credentials")
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

func (s *AuthService) ValidateAdminAPIKey(apiKey string) error {
	expectedKey := s.config.AdminAPIKey
	if subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedKey)) != 1 {
		return ErrInvalidAdminCredentials
	}

	return nil
}
