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
	sanitizeLocationData(graph)

	return graph, nil
}

// sanitizeLocationData strips exact coordinates from location gate nodes
// so the client never learns the real lat/lng. Only the hint (if enabled)
// and unlock radius are sent.
func sanitizeLocationData(graph *models.Graph) {
	if graph.Nodes == nil {
		return
	}
	for i := range *graph.Nodes {
		sanitizeNodeLocationData(&(*graph.Nodes)[i])
	}
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

		// Auto-unlock any reward nodes that became accessible. Reward nodes
		// have no gate mechanism so they should unlock automatically. Loop
		// because unlocking a reward may reveal further nodes one hop away.
		// The final iteration (where nothing was unlocked) gives us graphAfter.
		var graphAfter *models.Graph
		for {
			graph, err := repo.GetAccessibleGraph(ctx, node.GraphID)
			if err != nil {
				return fmt.Errorf("could not get accessible graph during reward auto-unlock: %w", err)
			}

			unlocked := false
			for _, n := range *graph.Nodes {
				if n.Type == models.RewardNode && n.UnlockedAt == nil {
					if err := repo.UnlockNode(ctx, n.ID); err != nil {
						return fmt.Errorf("failed to auto-unlock reward node %s: %w", n.ID, err)
					}
					unlocked = true
				}
			}

			if !unlocked {
				graphAfter = graph
				break
			}
		}

		// Get the unlocked node (fresh from DB with updated timestamp)
		unlockedNode, err := repo.GetAccessibleNode(ctx, nodeID)
		if err != nil {
			return fmt.Errorf("failed to get unlocked node: %w", err)
		}

		// Compute diff by subtracting pre-existing nodes/edges by ID
		newNodes := setDifferenceBy(*graphAfter.Nodes, *graphBefore.Nodes, func(n models.Node) string { return n.ID })
		newEdges := setDifferenceBy(*graphAfter.Edges, *graphBefore.Edges, func(e models.Edge) string { return e.ID })

		// Grant favour points from any newly revealed reward nodes.
		var totalFavours int
		for _, n := range newNodes {
			if n.Type == models.RewardNode {
				if rd, ok := n.Data.(*models.RewardData); ok && rd.GiveFavours > 0 {
					totalFavours += rd.GiveFavours
				}
			}
		}
		if totalFavours > 0 {
			count, err := repo.GetFavourCount(ctx)
			if err != nil {
				return fmt.Errorf("failed to get favour count: %w", err)
			}
			if err := repo.UpdateFavourCount(ctx, count.Total+totalFavours); err != nil {
				return fmt.Errorf("failed to update favour count: %w", err)
			}
		}

		result = &models.UnlockResult{
			UnlockedNode: *unlockedNode,
			NewNodes:     newNodes,
			NewEdges:     newEdges,
		}

		return nil
	})

	if err == nil && result != nil {
		sanitizeNodeLocationData(&result.UnlockedNode)
		for i := range result.NewNodes {
			sanitizeNodeLocationData(&result.NewNodes[i])
		}
	}

	return result, err
}

func sanitizeNodeLocationData(node *models.Node) {
	if node.Type != models.LocationGateNode {
		return
	}
	locData, ok := node.Data.(*models.LocationData)
	if !ok || locData == nil {
		return
	}
	node.Data = &models.ClientLocationData{
		RadiusM: locData.RadiusM,
		Hint:    locData.Hint,
	}
}

func validateUnlockPayload(node *models.Node, payload string) error {
	switch node.Type {
	case models.CodeGateNode:
		return validateCodeGatePayload(node, payload)
	case models.LocationGateNode:
		return validateLocationGatePayload(node, payload)
	case models.ManualNode:
		return validateManualGateUnlock(node)
	case models.TimeGateNode:
		return validateTimeGateUnlock(node)
	default:
		return fmt.Errorf("unsupported node type %s for unlocking", node.Type)
	}
}

func validateManualGateUnlock(node *models.Node) error {
	nodeData, ok := node.Data.(*models.ManualData)
	if !ok {
		return fmt.Errorf("invalid node data for manual gate node %s", node.ID)
	}

	if nodeData.UnlockedAt == nil {
		return NodeUnlockIncorrect
	}

	return nil
}

func validateTimeGateUnlock(node *models.Node) error {
	nodeData, ok := node.Data.(*models.TimeData)
	if !ok {
		return fmt.Errorf("invalid node data for time gate node %s", node.ID)
	}

	if time.Now().UTC().Before(nodeData.UnlockAt) {
		return NodeUnlockIncorrect
	}

	return nil
}

func validateCodeGatePayload(node *models.Node, payload string) error {
	nodeData, ok := node.Data.(*models.CodeData)
	if !ok {
		return fmt.Errorf("invalid node data for code gate node %s", node.ID)
	}

	for _, code := range nodeData.Codes {
		if strings.EqualFold(code, payload) {
			return nil
		}
	}

	return NodeUnlockIncorrect
}

func validateLocationGatePayload(node *models.Node, payload string) error {
	nodeData, ok := node.Data.(*models.LocationData)
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
