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
	FavourRepository
	GraphRepository
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

type FavourRepository interface {
	CreateFavourChoice(ctx context.Context, choice *models.FavourChoice) error // The ID of the newly created choice is set in the provided choice object.
	UpdateFavourChoice(ctx context.Context, choice models.FavourChoice) error
	DeleteFavourChoice(ctx context.Context, choiceID string) error
	ListFavourChoices(ctx context.Context) ([]models.FavourChoice, error)
	GetFavourCount(ctx context.Context) (models.FavourCount, error)
	UpdateFavourCount(ctx context.Context, total int) error
	CreateFavourRequest(ctx context.Context, request *models.FavourRequest) error // The ID of the newly created request is set in the provided request object.
	ListFavourRequests(ctx context.Context) ([]models.FavourRequest, error)
	DeleteFavourRequest(ctx context.Context, requestID string) error
	FulfilFavourRequest(ctx context.Context, requestID string) error
	UnfulfilFavourRequest(ctx context.Context, requestID string) error
}

type GraphRepository interface {
	ListStartNodes(ctx context.Context) ([]models.Node, error)
	DeleteGraph(ctx context.Context, startNodeID string) error
	CreateGraph(ctx context.Context, req models.NewGraphRequest) (string, error)
	GetGraph(ctx context.Context, startNodeID string) (*models.Graph, error)
}
