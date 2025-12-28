package repository

import (
	"backend/internal/models"
	"context"
)

type Repository interface {
	Close()
	WithTx(ctx context.Context, fn func(Repository) error) error

	RegistrationRepository
}

type RegistrationRepository interface {
	RefreshRegistrationCode(ctx context.Context, code models.RegistrationCode) error
	DeleteRegistrationCode(ctx context.Context) error
	GetRegistrationCode(ctx context.Context) (*models.RegistrationCode, error)
}
