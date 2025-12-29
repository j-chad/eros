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
