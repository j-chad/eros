package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type FavourService struct {
	repo repository.Repository
}

func NewFavourService(repo repository.Repository) *FavourService {
	return &FavourService{repo: repo}
}

func (s *FavourService) ListFavourChoices(ctx context.Context) ([]models.FavourChoice, error) {
	return s.repo.ListFavourChoices(ctx)
}

func (s *FavourService) GetFavourCount(ctx context.Context) (models.FavourCount, error) {
	return s.repo.GetFavourCount(ctx)
}

func (s *FavourService) ListFavourRequests(ctx context.Context) ([]models.FavourRequest, error) {
	return s.repo.ListFavourRequests(ctx)
}

func (s *FavourService) RequestFavour(ctx context.Context, request *models.FavourRequest) error {
	return s.repo.CreateFavourRequest(ctx, request)
}

func (s *FavourService) DeleteFavourRequest(ctx context.Context, requestID string) error {
	return s.repo.DeleteFavourRequest(ctx, requestID)
}
