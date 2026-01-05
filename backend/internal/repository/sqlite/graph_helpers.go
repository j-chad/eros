package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

// scanNodeFull scans a row from node_full view into a Node struct
func (s *sqliteDB) scanNodeFull(scanner interface {
	Scan(dest ...interface{}) error
}) (models.Node, error) {
	var node models.Node
	var uiPositionX sql.NullFloat64
	var uiPositionY sql.NullFloat64
	var dataJSON []byte

	err := scanner.Scan(
		&node.ID,
		&node.Type,
		&node.Title,
		&node.Description,
		&uiPositionX,
		&uiPositionY,
		&node.CreatedAt,
		&node.UpdatedAt,
		&node.UnlockedAt,

		&dataJSON,
	)
	if err != nil {
		return node, fmt.Errorf("failed to scan node: %w", err)
	}

	if uiPositionX.Valid && uiPositionY.Valid {
		node.UIPosition = &models.NodePosition{
			X: uiPositionX.Float64,
			Y: uiPositionY.Float64,
		}
	}

	// Populate type-specific data based on node type
	switch node.Type {
	case models.StartNode:
		node.Data = &models.StartData{}
		if err := json.Unmarshal(dataJSON, &node.Data); err != nil {
			return node, fmt.Errorf("failed to unmarshal start data: %w", err)
		}
	case models.LocationGateNode:
		node.Data = &models.LocationData{}
		if err := json.Unmarshal(dataJSON, &node.Data); err != nil {
			return node, fmt.Errorf("failed to unmarshal location data: %w", err)
		}
	case models.CodeGateNode:
		node.Data = &models.CodeData{}
		if err := json.Unmarshal(dataJSON, &node.Data); err != nil {
			return node, fmt.Errorf("failed to unmarshal code data: %w", err)
		}
	case models.RewardNode:
		node.Data = &models.RewardData{}
		if err := json.Unmarshal(dataJSON, &node.Data); err != nil {
			return node, fmt.Errorf("failed to unmarshal reward data: %w", err)
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

func (s *sqliteDB) updateNode(ctx context.Context, node models.Node) error {
	return s.withTx(ctx, func(txRepo *sqliteDB) error {
		// update node data
		var result sql.Result
		var err error
		switch node.Type {
		case models.StartNode:
			data, ok := node.Data.(*models.StartData)
			if !ok || data == nil {
				return fmt.Errorf("invalid start data")
			}

			result, err = s.executor().ExecContext(ctx, `
				UPDATE node_start SET viewport_x = ?, viewport_y = ?, viewport_zoom = ?
				WHERE node_id = ?
			`, data.ViewportX, data.ViewportY, data.ViewportZoom, node.ID)
		case models.LocationGateNode:
			data, ok := node.Data.(*models.LocationData)
			if !ok || data == nil {
				return fmt.Errorf("invalid location data")
			}

			result, err = s.executor().ExecContext(ctx, `
				INSERT INTO node_location_gate (node_id, latitude, longitude, radius_meters) VALUES (?, ?, ?, ?)
			    ON CONFLICT DO UPDATE SET latitude = ?2, longitude = ?3, radius_meters = ?4
			`, node.ID, data.Latitude, data.Longitude, data.RadiusM)
		default:
			return fmt.Errorf("unknown node type: %s", node.Type)
		}

		if err != nil {
			return fmt.Errorf("failed to update node data: %w", err)
		}

		affected, err := result.RowsAffected()
		if err != nil || affected == 0 {
			return fmt.Errorf("no node data updated: %w", err)
		}

		return nil
	})
}
