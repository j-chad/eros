package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

func (s *sqliteDB) upsertNodeData(ctx context.Context, nodeID string, nodeType models.NodeType, data any) error {
	if data == nil {
		return nil
	}

	switch nodeType {
	case models.ManualNode:
		manualData, ok := data.(models.ManualData)
		if !ok {
			return fmt.Errorf("invalid data")
		}
		_, err := s.executor().ExecContext(ctx, `
			INSERT INTO node_manual_gate (
			    node_id, instructions, unlocked_at
			) VALUES (?, ?, ?)
			ON CONFLICT (node_id) DO UPDATE SET
				instructions = excluded.instructions,
				unlocked_at = excluded.unlocked_at
		`, nodeID, manualData.Instructions, manualData.UnlockedAt)
		return err
	case models.LocationGateNode:
		locationData, ok := data.(models.LocationData)
		if !ok {
			return fmt.Errorf("invalid data type for location node")
		}
		var hintLat, hintLng *float64
		var hintRadius *int
		if locationData.Hint != nil {
			hintLat = &locationData.Hint.Latitude
			hintLng = &locationData.Hint.Longitude
			hintRadius = &locationData.Hint.RadiusM
		}
		_, err := s.executor().ExecContext(ctx, `
			INSERT INTO node_location_gate (
				node_id,
				latitude,
				longitude,
				radius_meters,
				show_hint,
				hint_latitude,
				hint_longitude,
				hint_radius_meters
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (node_id) DO UPDATE SET
				latitude = excluded.latitude,
				longitude = excluded.longitude,
				radius_meters = excluded.radius_meters,
				show_hint = excluded.show_hint,
				hint_latitude = excluded.hint_latitude,
				hint_longitude = excluded.hint_longitude,
				hint_radius_meters = excluded.hint_radius_meters
		`, nodeID, locationData.Latitude, locationData.Longitude, locationData.RadiusM,
			locationData.ShowHint, hintLat, hintLng, hintRadius)
		return err
	case models.CodeGateNode:
		codeData, ok := data.(models.CodeData)
		if !ok {
			return fmt.Errorf("invalid data type for code node")
		}
		codes := codeData.Codes
		if codes == nil {
			codes = []string{}
		}
		codesJSON, err := json.Marshal(codes)
		if err != nil {
			return fmt.Errorf("failed to marshal codes: %w", err)
		}
		_, err = s.executor().ExecContext(ctx, `
			INSERT INTO node_code_gate (
				node_id,
				codes
			) VALUES (?, ?)
			ON CONFLICT (node_id) DO UPDATE SET
				codes = excluded.codes
		`, nodeID, string(codesJSON))
		return err
	case models.TimeGateNode:
		timeData, ok := data.(models.TimeData)
		if !ok {
			return fmt.Errorf("invalid data type for time node")
		}
		_, err := s.executor().ExecContext(ctx, `
			INSERT INTO node_time_gate (
				node_id,
				unlock_at
			) VALUES (?, ?)
			ON CONFLICT (node_id) DO UPDATE SET
				unlock_at = excluded.unlock_at
		`, nodeID, timeData.UnlockAt)
		return err
	case models.RewardNode:
		rewardData, ok := data.(models.RewardData)
		if !ok {
			return fmt.Errorf("invalid data type for reward node")
		}
		_, err := s.executor().ExecContext(ctx, `
			INSERT INTO node_reward (
				node_id,
				reward_type,
			    payload,
				give_favours
			) VALUES (?, ?, ?, ?)
			ON CONFLICT (node_id) DO UPDATE SET
				reward_type = excluded.reward_type,
				payload = excluded.payload,
				give_favours = excluded.give_favours
		`, nodeID, rewardData.RewardType, rewardData.Payload, rewardData.GiveFavours)
		return err
	case models.StartNode:
		// No additional data to insert
		return nil
	default:
		return fmt.Errorf("unknown node type: %s", nodeType)
	}
}

func (s *sqliteDB) createNode(ctx context.Context, graphID string, node models.Node) error {
	return s.withTx(ctx, nil, func(tx *sqliteDB) error {
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

		_, err := tx.executor().ExecContext(ctx, `
			INSERT INTO node (
			    id,
				graph_id,
				type,
				title,
				description,
				ui_pos_x,
				ui_pos_y
		  	) VALUES (?, ?, ?, ?, ?, ?, ?)
		`, node.ID, graphID, node.Type, node.Title, node.Description, uiPositionX, uiPositionY)
		if err != nil {
			return fmt.Errorf("failed to insert node: %w", err)
		}

		// insert type-specific data
		err = tx.upsertNodeData(ctx, node.ID, node.Type, node.Data)
		if err != nil {
			return fmt.Errorf("failed to insert node data: %w", err)
		}

		return nil
	})
}

func (s *sqliteDB) createEdge(ctx context.Context, graphID string, edge models.Edge) error {
	_, err := s.executor().ExecContext(ctx, `
		INSERT INTO edge (
		  	id,
			graph_id,
			source_node_id,
			destination_node_id,
			choice_label
		) VALUES (?, ?, ?, ?, ?)
	`, edge.ID, graphID, edge.From, edge.To, edge.ChoiceLabel)
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

	return s.withTx(ctx, nil, func(tx *sqliteDB) error {
		for _, node := range nodes {
			// since the node still exists, remove it from the deleted set
			delete(deletedNodeIDSet, node.ID)

			_, exists := existingNodeIDSet[node.ID]
			if node.ID != "" && exists {
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

	return s.withTx(ctx, nil, func(tx *sqliteDB) error {
		for _, edge := range edges {
			// since the edge still exists, remove it from the deleted set
			delete(deletedEdgeIDSet, edge.ID)

			_, exists := existingEdgeIDSet[edge.ID]
			if edge.ID != "" && exists {
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
