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

func (s *GraphService) ListGraphs(ctx context.Context) ([]models.Graph, error) {
	graphs, err := s.repo.ListGraphs(ctx)
	if err != nil {
		return nil, err
	}

	sanitisedGraphs := make([]models.Graph, 0, len(graphs))
	for _, graph := range graphs {
		if graph.StartingAt.Before(time.Now()) {
			sanitisedGraphs = append(sanitisedGraphs, models.Graph{
				ID:          graph.ID,
				Title:       graph.Title,
				Description: graph.Description,
				StartingAt:  graph.StartingAt,
			})
			continue
		}

		// For graphs that haven't started yet, we only return the ID and starting time to avoid leaking information about them.
		sanitisedGraphs = append(sanitisedGraphs, models.Graph{
			ID:         graph.ID,
			StartingAt: graph.StartingAt,
		})
	}

	return sanitisedGraphs, nil
}

func (s *GraphService) GetGraph(ctx context.Context, graphID string) (*models.Graph, error) {
	graph, err := s.repo.GetAccessibleGraph(ctx, graphID)
	if err != nil {
		return nil, err
	}

	return graph, nil
}
