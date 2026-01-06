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
		&node.GraphID,
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
	case models.StartNode, models.ChoiceNode:
		// No additional data for these types
	default:
		return node, fmt.Errorf("unknown node type: %s", node.Type)
	}

	return node, nil
}

// getCompleteNodes retrieves all nodes reachable from a start node with full data
func (s *sqliteDB) getCompleteNodes(ctx context.Context, graphID string) ([]models.Node, error) {
	rows, err := s.executor().QueryContext(ctx, `
        SELECT * FROM node_full
        WHERE graph_id = ?
    `, graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to query graph nodes: %w", err)
	}
	defer rows.Close()

	nodes := make([]models.Node, 0)

	for rows.Next() {
		node, err := s.scanNodeFull(rows)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating nodes: %w", err)
	}

	return nodes, nil
}

// getEdges retrieves all edges for the given node IDs
func (s *sqliteDB) getEdges(ctx context.Context, graphID string) ([]models.Edge, error) {
	rows, err := s.executor().QueryContext(ctx, `
        SELECT 
            id,
            source_node_id,
            destination_node_id,
            choice_label,
            created_at,
            updated_at
        FROM edge
        WHERE graph_id = ?
    `, graphID)
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

func (s *sqliteDB) getNodeIDs(ctx context.Context, graphID string) ([]string, error) {
	rows, err := s.executor().QueryContext(ctx, `
		SELECT id FROM node
		WHERE graph_id = ?
	`, graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to query node IDs: %w", err)
	}
	defer rows.Close()

	nodeIDs := make([]string, 0)
	for rows.Next() {
		var nodeID string
		if err := rows.Scan(&nodeID); err != nil {
			return nil, fmt.Errorf("failed to scan node ID: %w", err)
		}
		nodeIDs = append(nodeIDs, nodeID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating node IDs: %w", err)
	}

	return nodeIDs, nil
}

func (s *sqliteDB) getEdgeIDs(ctx context.Context, graphID string) ([]string, error) {
	rows, err := s.executor().QueryContext(ctx, `
		SELECT id FROM edge
		WHERE graph_id = ?
	`, graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to query edge IDs: %w", err)
	}
	defer rows.Close()

	edgeIDs := make([]string, 0)
	for rows.Next() {
		var edgeID string
		if err := rows.Scan(&edgeID); err != nil {
			return nil, fmt.Errorf("failed to scan edge ID: %w", err)
		}
		edgeIDs = append(edgeIDs, edgeID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating edge IDs: %w", err)
	}

	return edgeIDs, nil
}

func (s *sqliteDB) upsertNodeData(ctx context.Context, nodeID string, nodeType models.NodeType, data any) error {
	switch nodeType {
	case models.LocationGateNode:
		locationData, ok := data.(*models.LocationData)
		if !ok {
			return fmt.Errorf("invalid data type for location node")
		}
		_, err := s.executor().ExecContext(ctx, `
			INSERT INTO node_location_gate (
				node_id,
				latitude,
				longitude,
				radius_meters
			) VALUES (?, ?, ?, ?)
			ON CONFLICT (node_id) DO UPDATE SET
				latitude = excluded.latitude,
				longitude = excluded.longitude,
				radius_meters = excluded.radius_meters
		`, nodeID, locationData.Latitude, locationData.Longitude, locationData.RadiusM)
		return err
	case models.CodeGateNode:
		codeData, ok := data.(*models.CodeData)
		if !ok {
			return fmt.Errorf("invalid data type for code node")
		}
		_, err := s.executor().ExecContext(ctx, `
			INSERT INTO node_code_gate (
				node_id,
				code
			) VALUES (?, ?)
			ON CONFLICT (node_id) DO UPDATE SET
				code = excluded.code
		`, nodeID, codeData.Code)
		return err
	case models.RewardNode:
		rewardData, ok := data.(*models.RewardData)
		if !ok {
			return fmt.Errorf("invalid data type for reward node")
		}
		_, err := s.executor().ExecContext(ctx, `
			INSERT INTO node_reward (
				node_id,
				reward_type,
			    content_html,
			    content_media_type,
				give_favours
			) VALUES (?, ?, ?, ?, ?)
			ON CONFLICT (node_id) DO UPDATE SET
				reward_type = excluded.reward_type,
				content_html = excluded.content_html,
				content_media_type = excluded.content_media_type,
				give_favours = excluded.give_favours
		`, nodeID, rewardData.RewardType, rewardData.Content, rewardData.MediaType, rewardData.GiveFavours)
		return err
	case models.StartNode, models.ChoiceNode:
		// No additional data to insert
		return nil
	default:
		return fmt.Errorf("unknown node type: %s", nodeType)
	}
}

func (s *sqliteDB) createNode(ctx context.Context, graphID string, node models.Node) error {
	return s.withTx(ctx, func(tx *sqliteDB) error {
		// insert node
		var uiPositionX sql.NullFloat64
		var uiPositionY sql.NullFloat64
		if node.UIPosition != nil {
			uiPositionX = sql.NullFloat64{Float64: node.UIPosition.X, Valid: true}
			uiPositionY = sql.NullFloat64{Float64: node.UIPosition.Y, Valid: true}
		} else {
			uiPositionX = sql.NullFloat64{Valid: false}
			uiPositionY = sql.NullFloat64{Valid: false}
		}

		result, err := tx.executor().ExecContext(ctx, `
			INSERT INTO node (
				graph_id,
				type,
				title,
				description,
				ui_pos_x,
				ui_pos_y  
		  	) VALUES (?, ?, ?, ?, ?, ?)
		`, graphID, node.Type, node.Title, node.Description, uiPositionX, uiPositionY)
		if err != nil {
			return fmt.Errorf("failed to insert node: %w", err)
		}

		nodeID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for node: %w", err)
		}

		// insert type-specific data
		err = tx.upsertNodeData(ctx, strconv.FormatInt(nodeID, 10), node.Type, node.Data)
		if err != nil {
			return fmt.Errorf("failed to insert node data: %w", err)
		}

		return nil
	})
}

func (s *sqliteDB) createEdge(ctx context.Context, graphID string, edge models.Edge) error {
	_, err := s.executor().ExecContext(ctx, `
		INSERT INTO edge (
			graph_id,
			source_node_id,
			destination_node_id,
			choice_label
		) VALUES (?, ?, ?, ?)
	`, graphID, edge.From, edge.To, edge.ChoiceLabel)
	if err != nil {
		return fmt.Errorf("failed to insert edge: %w", err)
	}

	return nil
}

func (s *sqliteDB) updateNode(ctx context.Context, node models.Node) error {
	// update node
	var uiPositionX sql.NullFloat64
	var uiPositionY sql.NullFloat64
	if node.UIPosition != nil {
		uiPositionX = sql.NullFloat64{Float64: node.UIPosition.X, Valid: true}
		uiPositionY = sql.NullFloat64{Float64: node.UIPosition.Y, Valid: true}
	} else {
		uiPositionX = sql.NullFloat64{Valid: false}
		uiPositionY = sql.NullFloat64{Valid: false}
	}

	_, err := s.executor().ExecContext(ctx, `
		UPDATE node SET
			type = ?,
			title = ?,
			description = ?,
			ui_pos_x = ?,
			ui_pos_y = ?
		WHERE id = ?
	`, node.Type, node.Title, node.Description, uiPositionX, uiPositionY, node.ID)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}

	// update type-specific data
	err = s.upsertNodeData(ctx, node.ID, node.Type, node.Data)
	if err != nil {
		return fmt.Errorf("failed to update node data: %w", err)
	}

	return nil
}

func (s *sqliteDB) updateEdge(ctx context.Context, edge models.Edge) error {
	_, err := s.executor().ExecContext(ctx, `
		UPDATE edge SET
			source_node_id = ?,
			destination_node_id = ?,
			choice_label = ?
		WHERE id = ?
	`, edge.From, edge.To, edge.ChoiceLabel, edge.ID)
	if err != nil {
		return fmt.Errorf("failed to update edge: %w", err)
	}

	return nil
}

func (s *sqliteDB) deleteNodes(ctx context.Context, nodeIDs []string) error {
	jsonIDs, err := json.Marshal(nodeIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal node IDs: %w", err)
	}

	_, err = s.executor().ExecContext(ctx, `
		DELETE FROM node
		WHERE id IN (SELECT value FROM json_each(?))
	`, string(jsonIDs))

	if err != nil {
		return fmt.Errorf("failed to delete nodes: %w", err)
	}

	return nil
}

func (s *sqliteDB) deleteEdges(ctx context.Context, edgeIDs []string) error {
	jsonIDs, err := json.Marshal(edgeIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal edge IDs: %w", err)
	}

	_, err = s.executor().ExecContext(ctx, `
		DELETE FROM edge
		WHERE id IN (SELECT value FROM json_each(?))
	`, string(jsonIDs))

	if err != nil {
		return fmt.Errorf("failed to delete edges: %w", err)
	}

	return nil
}

func (s *sqliteDB) updateNodes(ctx context.Context, graphID string, nodes []models.Node) error {
	existingNodeIDs, err := s.getNodeIDs(ctx, graphID)
	if err != nil {
		return fmt.Errorf("failed to get existing node IDs: %w", err)
	}

	existingNodeIDSet := make(map[string]struct{}, len(existingNodeIDs))
	deletedNodeIDSet := make(map[string]struct{}, len(existingNodeIDs))
	for _, id := range existingNodeIDs {
		existingNodeIDSet[id] = struct{}{}
		deletedNodeIDSet[id] = struct{}{}
	}

	return s.withTx(ctx, func(tx *sqliteDB) error {
		for _, node := range nodes {
			// since the node still exists, remove it from the deleted set
			delete(deletedNodeIDSet, node.ID)

			if node.ID != "" && existingNodeIDSet[node.ID] == struct{}{} {
				err := tx.updateNode(ctx, node)
				if err != nil {
					return fmt.Errorf("failed to update node %s: %w", node.ID, err)
				}
			} else {
				err := tx.createNode(ctx, graphID, node)
				if err != nil {
					return fmt.Errorf("failed to create node %s: %w", node.ID, err)
				}
			}
		}

		if len(deletedNodeIDSet) > 0 {
			idsToDelete := make([]string, 0, len(deletedNodeIDSet))
			for id := range deletedNodeIDSet {
				idsToDelete = append(idsToDelete, id)
			}

			err := tx.deleteNodes(ctx, idsToDelete)
			if err != nil {
				return fmt.Errorf("failed to delete nodes: %w", err)
			}
		}

		return nil
	})
}

func (s *sqliteDB) updateEdges(ctx context.Context, graphID string, edges []models.Edge) error {
	existingEdgeIDs, err := s.getEdgeIDs(ctx, graphID)
	if err != nil {
		return fmt.Errorf("failed to get existing edge IDs: %w", err)
	}

	existingEdgeIDSet := make(map[string]struct{}, len(existingEdgeIDs))
	deletedEdgeIDSet := make(map[string]struct{}, len(existingEdgeIDs))
	for _, id := range existingEdgeIDs {
		existingEdgeIDSet[id] = struct{}{}
		deletedEdgeIDSet[id] = struct{}{}
	}

	return s.withTx(ctx, func(tx *sqliteDB) error {
		for _, edge := range edges {
			// since the edge still exists, remove it from the deleted set
			delete(deletedEdgeIDSet, edge.ID)

			if edge.ID != "" && existingEdgeIDSet[edge.ID] == struct{}{} {
				err := tx.updateEdge(ctx, edge)
				if err != nil {
					return fmt.Errorf("failed to update edge %s: %w", edge.ID, err)
				}
			} else {
				err := tx.createEdge(ctx, graphID, edge)
				if err != nil {
					return fmt.Errorf("failed to create edge %s: %w", edge.ID, err)
				}
			}
		}

		if len(deletedEdgeIDSet) > 0 {
			idsToDelete := make([]string, 0, len(deletedEdgeIDSet))
			for id := range deletedEdgeIDSet {
				idsToDelete = append(idsToDelete, id)
			}

			err := tx.deleteEdges(ctx, idsToDelete)
			if err != nil {
				return fmt.Errorf("failed to delete edges: %w", err)
			}
		}

		return nil
	})
}
