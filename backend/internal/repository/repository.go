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
	FavourChoiceRepository
}

type RegistrationRepository interface {
	RefreshRegistrationCode(ctx context.Context, code models.RegistrationCode) error
	DeleteRegistrationCode(ctx context.Context) error
	GetRegistrationCode(ctx context.Context) (*models.RegistrationCode, error)
}

type DeviceRepository interface {
	RegisterDevice(ctx context.Context, token string, deviceInfo string, expiry time.Time) error
	GetDeviceExpiryByToken(ctx context.Context, token string) (expiresAt *time.Time, err error)
	ListDevices(ctx context.Context) ([]models.Device, error)
	DeleteDevice(ctx context.Context, deviceID string) error
	UpdateDeviceInfo(ctx context.Context, deviceID string, deviceInfo string) error
	UpdateDeviceLastSeenByToken(ctx context.Context, token string, lastSeen time.Time) error
}

type FavourChoiceRepository interface {
	CreateFavourChoice(ctx context.Context, choice *models.FavourChoice) error // The ID of the newly created choice is set in the provided choice object.
	UpdateFavourChoice(ctx context.Context, choice models.FavourChoice) error
	DeleteFavourChoice(ctx context.Context, choiceID string) error
	ListFavourChoices(ctx context.Context) ([]models.FavourChoice, error)
	GetFavourCount(ctx context.Context) (int, error)
	UpdateFavourCount(ctx context.Context, count int) error
}
