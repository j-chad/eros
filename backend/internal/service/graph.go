package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type GraphService struct {
	repo repository.Repository
}

var NodeUnlockIncorrect = errors.New("node unlock incorrect")

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
				Status:      graph.Status,
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

	graph.Viewport = nil

	return graph, nil
}

func (s *GraphService) UnlockNode(ctx context.Context, nodeID string, payload string) (*models.UnlockResult, error) {
	node, err := s.repo.GetAccessibleNode(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	if err := validateUnlockPayload(node, payload); err != nil {
		return nil, err
	}

	var result *models.UnlockResult

	err = s.repo.WithTx(ctx, nil, func(repo repository.Repository) error {
		graphBefore, err := repo.GetAccessibleGraph(ctx, node.GraphID)
		if err != nil {
			return fmt.Errorf("could not get accessible before unlock: %w", err)
		}

		if err := repo.UnlockNode(ctx, nodeID); err != nil {
			return fmt.Errorf("failed to unlock node: %w", err)
		}

		graphAfter, err := repo.GetAccessibleGraph(ctx, node.GraphID)
		if err != nil {
			return fmt.Errorf("could not get accessible after unlock: %w", err)
		}

		// Get the unlocked node (fresh from DB with updated timestamp)
		unlockedNode, err := repo.GetAccessibleNode(ctx, nodeID)
		if err != nil {
			return fmt.Errorf("failed to get unlocked node: %w", err)
		}

		// Compute diff by subtracting pre-existing nodes/edges from the new set
		newNodes := setDifference(*graphAfter.Nodes, *graphBefore.Nodes)
		newEdges := setDifference(*graphAfter.Edges, *graphBefore.Edges)

		result = &models.UnlockResult{
			UnlockedNode: *unlockedNode,
			NewNodes:     newNodes,
			NewEdges:     newEdges,
		}

		return nil
	})

	return result, err
}

func validateUnlockPayload(node *models.Node, payload string) error {
	switch node.Type {
	case models.CodeGateNode:
		return validateCodeGatePayload(node, payload)
	case models.LocationGateNode:
		return validateLocationGatePayload(node, payload)
	case models.ManualNode:
		return validateManualGateUnlock(node)
	default:
		return fmt.Errorf("unsupported node type %s for unlocking", node.Type)
	}
}

func validateManualGateUnlock(node *models.Node) error {
	nodeData, ok := node.Data.(models.ManualData)
	if !ok {
		return fmt.Errorf("invalid node data for manual gate node %s", node.ID)
	}

	if nodeData.UnlockedAt == nil {
		return NodeUnlockIncorrect
	}

	return nil
}

func validateCodeGatePayload(node *models.Node, payload string) error {
	nodeData, ok := node.Data.(models.CodeData)
	if !ok {
		return fmt.Errorf("invalid node data for code gate node %s", node.ID)
	}

	if !strings.EqualFold(nodeData.Code, payload) {
		return NodeUnlockIncorrect
	}

	return nil
}

func validateLocationGatePayload(node *models.Node, payload string) error {
	nodeData, ok := node.Data.(models.LocationData)
	if !ok {
		return fmt.Errorf("invalid node data for location gate node %s", node.ID)
	}

	var lat, lng float64
	_, err := fmt.Sscanf(payload, "%f,%f", &lat, &lng)
	if err != nil {
		return fmt.Errorf("invalid payload format for location gate node %s: %w", node.ID, err)
	}

	radius := nodeData.RadiusM
	if radius <= 0 {
		radius = 10
	}

	distance := haversineDistance(lat, lng, nodeData.Latitude, nodeData.Longitude)
	if distance > float64(radius) {
		return NodeUnlockIncorrect
	}

	return nil
}
