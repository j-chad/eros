package repository

import (
	"backend/internal/models"
	"context"
	"time"
)

type Repository interface {
	Close()
	WithTx(ctx context.Context, fn func(Repository) error) error

	RegistrationRepository
	DeviceRepository
}

type RegistrationRepository interface {
	RefreshRegistrationCode(ctx context.Context, code models.RegistrationCode) error
	DeleteRegistrationCode(ctx context.Context) error
	GetRegistrationCode(ctx context.Context) (*models.RegistrationCode, error)
}

type DeviceRepository interface {
	RegisterDevice(ctx context.Context, token string, deviceInfo string, expiry time.Time) error
	GetDeviceExpiryByToken(ctx context.Context, token string) (expiresAt *time.Time, err error)
}
