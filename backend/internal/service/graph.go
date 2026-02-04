package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"time"
)

type GraphService struct {
	repo repository.Repository
}

func NewGraphService(repo repository.Repository) *GraphService {
	return &GraphService{repo: repo}
}

func (s *GraphService) GetGraph(ctx context.Context, graphID string) (*models.Graph, error) {
	graph, err := s.repo.GetGraph(ctx, graphID)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	if graph.StartingAt.After(currentTime) {
		return nil, nil
	}

	// TODO: shake graph down to only include available nodes and edges
	// shakeGraph(graph)

	return graph, nil
}

func shakeGraph(graph *models.Graph) {
}
