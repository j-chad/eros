package service

import (
	"backend/internal/crypto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"time"
)

const RegistrationExpiryDuration = 10 * time.Minute

type AdminService struct {
	repo repository.Repository
}

func NewAdminService(repo repository.Repository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) CreateRegistrationCode(ctx context.Context) (models.RegistrationCode, error) {
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

func (s *AdminService) InvalidateRegistrationCode(ctx context.Context) error {
	return s.repo.DeleteRegistrationCode(ctx)
}

func (s *AdminService) GetRegistrationCode(ctx context.Context) (*models.RegistrationCode, error) {
	return s.repo.GetRegistrationCode(ctx)
}

func (s *AdminService) ListDevices(ctx context.Context) ([]models.Device, error) {
	return s.repo.ListDevices(ctx)
}

func (s *AdminService) RevokeDevice(ctx context.Context, deviceID string) error {
	return s.repo.DeleteDevice(ctx, deviceID)
}

func (s *AdminService) UpdateDeviceInfo(ctx context.Context, deviceID string, deviceInfo string) error {
	return s.repo.UpdateDeviceInfo(ctx, deviceID, deviceInfo)
}

func (s *AdminService) CreateFavourChoice(ctx context.Context, choice *models.FavourChoice) error {
	return s.repo.CreateFavourChoice(ctx, choice)
}

func (s *AdminService) UpdateFavourChoice(ctx context.Context, choice models.FavourChoice) error {
	return s.repo.UpdateFavourChoice(ctx, choice)
}

func (s *AdminService) DeleteFavourChoice(ctx context.Context, choiceID string) error {
	return s.repo.DeleteFavourChoice(ctx, choiceID)
}

func (s *AdminService) UpdateFavourCount(ctx context.Context, count int) error {
	return s.repo.UpdateFavourCount(ctx, count)
}

func (s *AdminService) FulfilFavourRequest(ctx context.Context, requestID string) error {
	return s.repo.FulfilFavourRequest(ctx, requestID)
}
func (s *AdminService) UnfulfilFavourRequest(ctx context.Context, requestID string) error {
	return s.repo.UnfulfilFavourRequest(ctx, requestID)
}

func (s *AdminService) ListStartNodes(ctx context.Context) ([]models.Node, error) {
	return s.repo.ListStartNodes(ctx)
}

func (s *AdminService) CreateGraph(ctx context.Context, req models.NewGraphRequest) (string, error) {
	return s.repo.CreateGraph(ctx, req)
}

func (s *AdminService) DeleteGraph(ctx context.Context, startNodeID string) error {
	return s.repo.DeleteGraph(ctx, startNodeID)
}
