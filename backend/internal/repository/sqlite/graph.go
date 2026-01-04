package sqlite

import (
	"backend/internal/models"
	"context"
	"fmt"
)

func (s *sqliteDB) ListStartNodes(ctx context.Context) ([]models.Node, error) {
	query := `
        SELECT 
            n.id,
            n.type,
            n.title,
            n.description,
            n.created_at,
            n.updated_at,
            n.unlocked_at,
            ns.starting_at
        FROM node n
        INNER JOIN node_start ns ON n.id = ns.node_id
        ORDER BY n.id
    `

	rows, err := s.executor().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query start nodes: %w", err)
	}
	defer rows.Close()

	nodes := make([]models.Node, 0)

	for rows.Next() {
		var node models.Node
		var startData models.StartData

		err := rows.Scan(
			&node.ID,
			&node.Type,
			&node.Title,
			&node.Description,
			&node.CreatedAt,
			&node.UpdatedAt,
			&node.UnlockedAt,
			&startData.StartingAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan start node: %w", err)
		}

		node.Start = &startData

		nodes = append(nodes, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating start nodes: %w", err)
	}

	return nodes, nil
}

func (s *sqliteDB) DeleteGraph(ctx context.Context, startNodeID string) error {
	query := `
		DELETE FROM node 
		WHERE id IN (
			WITH RECURSIVE graph_nodes(node_id) AS (
				SELECT node_id 
				FROM node_start 
				WHERE node_id = ?
				
				UNION
				
				SELECT e.to_node_id
				FROM graph_nodes gn
				JOIN edge e ON gn.node_id = e.from_node_id
			)
			SELECT node_id FROM graph_nodes
		);
	`

	_, err := s.executor().ExecContext(ctx, query, startNodeID)
	if err != nil {
		return fmt.Errorf("failed to delete graph starting from node %s: %w", startNodeID, err)
	}

	return nil
}

func (s *sqliteDB) CreateGraph(ctx context.Context, req models.NewGraphRequest) (string, error) {
	var startNodeID int64

	err := s.withTx(ctx, func(txRepo *sqliteDB) error {
		// Insert start node
		result, err := txRepo.executor().ExecContext(ctx, `
			INSERT INTO node (type, title, description)
			VALUES (?, ?, ?)
		`, models.StartNode, req.Title, req.Description)
		if err != nil {
			return fmt.Errorf("failed to insert start node: %w", err)
		}

		startNodeID, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for start node: %w", err)
		}

		// Insert into node_start
		_, err = txRepo.executor().ExecContext(ctx, `
			INSERT INTO node_start (node_id, starting_at)
			VALUES (?, ?)
		`, startNodeID, req.StartingAt)
		if err != nil {
			return fmt.Errorf("failed to insert into node_start: %w", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", startNodeID), nil
}
