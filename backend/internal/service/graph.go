package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type GraphService struct {
	repo repository.Repository
}

func NewGraphService(repo repository.Repository) *GraphService {
	return &GraphService{repo: repo}
}

func (s *GraphService) GetGraph(ctx context.Context, graphID string) (*models.Graph, error) {
	graph, err := s.repo.GetAccessibleGraph(ctx, graphID)
	if err != nil {
		return nil, err
	}

	return graph, nil
}
