package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// scanNodeFull scans a row from node_full view into a Node struct
func (s *sqliteDB) scanNodeFull(scanner interface {
	Scan(dest ...interface{}) error
}) (models.Node, error) {
	var node models.Node
	var startingAt sql.NullTime
	var latitude, longitude sql.NullFloat64
	var radiusMeters sql.NullInt64
	var code sql.NullString
	var rewardType, contentHTML, contentMediaType sql.NullString
	var giveFavours sql.NullInt64

	err := scanner.Scan(
		&node.ID,
		&node.Type,
		&node.Title,
		&node.Description,
		&node.CreatedAt,
		&node.UpdatedAt,
		&node.UnlockedAt,
		&startingAt,
		&latitude,
		&longitude,
		&radiusMeters,
		&code,
		&rewardType,
		&contentHTML,
		&contentMediaType,
		&giveFavours,
	)
	if err != nil {
		return node, fmt.Errorf("failed to scan node: %w", err)
	}

	// Populate type-specific data based on node type
	switch node.Type {
	case models.StartNode:
		if !startingAt.Valid {
			return node, fmt.Errorf("node %d is type START but missing starting_at data", node.ID)
		}
		node.Start = &models.StartData{
			StartingAt: startingAt.Time.Format(time.RFC3339),
		}

	case models.LocationGateNode:
		if !latitude.Valid || !longitude.Valid || !radiusMeters.Valid {
			return node, fmt.Errorf("node %d is type LOCATION but missing location data", node.ID)
		}
		node.Location = &models.LocationData{
			Latitude:  latitude.Float64,
			Longitude: longitude.Float64,
			RadiusM:   int(radiusMeters.Int64),
		}

	case models.CodeGateNode:
		if !code.Valid {
			return node, fmt.Errorf("node %d is type CODE but missing code data", node.ID)
		}
		node.Code = &models.CodeData{
			Code: code.String,
		}

	case models.RewardNode:
		if !rewardType.Valid || !contentHTML.Valid || !contentMediaType.Valid || !giveFavours.Valid {
			return node, fmt.Errorf("node %d is type REWARD but missing reward data", node.ID)
		}
		node.Reward = &models.RewardData{
			RewardType:  rewardType.String,
			Content:     contentHTML.String,
			MediaType:   contentMediaType.String,
			GiveFavours: int(giveFavours.Int64),
		}

	case models.ChoiceNode:
		// Choice nodes don't have additional data (choices are in edges)
	default:
		return node, fmt.Errorf("unknown node type: %s", node.Type)
	}

	return node, nil
}

// getCompleteNodes retrieves all nodes reachable from a start node with full data
func (s *sqliteDB) getCompleteNodes(ctx context.Context, startNodeID string) (map[string]*models.Node, []int64, error) {
	query := `
        WITH RECURSIVE graph_nodes(node_id) AS (
            SELECT ? as node_id
            UNION
            SELECT e.to_node_id
            FROM graph_nodes gn
            JOIN edge e ON gn.node_id = e.from_node_id
        )
        SELECT * FROM node_full
        WHERE id IN (SELECT node_id FROM graph_nodes)
        ORDER BY id
    `

	rows, err := s.executor().QueryContext(ctx, query, startNodeID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query graph nodes: %w", err)
	}
	defer rows.Close()

	nodesMap := make(map[string]*models.Node)
	nodeIDs := make([]int64, 0)

	for rows.Next() {
		node, err := s.scanNodeFull(rows)
		if err != nil {
			return nil, nil, err
		}
		node.Edges = []models.Edge{}
		nodesMap[node.ID] = &node

		int64ID, err := strconv.ParseInt(node.ID, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert node ID to int64: %w", err)
		}
		nodeIDs = append(nodeIDs, int64ID)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating nodes: %w", err)
	}

	return nodesMap, nodeIDs, nil
}

// getCompleteEdges retrieves all edges for the given node IDs
func (s *sqliteDB) getCompleteEdges(ctx context.Context, nodeIDs []int64) ([]models.Edge, error) {
	if len(nodeIDs) == 0 {
		return []models.Edge{}, nil
	}

	nodeIDsJSON, err := json.Marshal(nodeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal node IDs: %w", err)
	}

	query := `
        SELECT 
            id,
            from_node_id,
            to_node_id,
            choice_label,
            created_at,
            updated_at
        FROM edge
        WHERE from_node_id IN (SELECT value FROM json_each(?))
        ORDER BY from_node_id, id
    `

	rows, err := s.executor().QueryContext(ctx, query, string(nodeIDsJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to query edges: %w", err)
	}
	defer rows.Close()

	edges := make([]models.Edge, 0)
	for rows.Next() {
		var edge models.Edge

		err := rows.Scan(
			&edge.ID,
			&edge.From,
			&edge.To,
			&edge.ChoiceLabel,
			&edge.CreatedAt,
			&edge.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan edge: %w", err)
		}

		edges = append(edges, edge)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating edges: %w", err)
	}

	return edges, nil
}
