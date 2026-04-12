package sqlite

import (
	"backend/internal/crypto"
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (s *sqliteDB) ListGraphs(ctx context.Context) ([]models.Graph, error) {
	query := `
        SELECT
            g.id,
            g.title,
            g.description,
            g.created_at,
            g.updated_at,
            g.starting_at,
            COUNT(n.id) AS total_nodes,
            COUNT(n.unlocked_at) AS unlocked_nodes,
            SUM(CASE WHEN n.type != 'start' AND n.unlocked_at IS NOT NULL THEN 1 ELSE 0 END) AS unlocked_non_start
        FROM graph g
        LEFT JOIN node n ON n.graph_id = g.id
        GROUP BY g.id
        ORDER BY g.starting_at DESC
    `

	rows, err := s.executor().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query graphs: %w", err)
	}
	defer rows.Close()

	graphs := make([]models.Graph, 0)

	for rows.Next() {
		var graph models.Graph
		var totalNodes, unlockedNodes, unlockedNonStart int

		err := rows.Scan(
			&graph.ID,
			&graph.Title,
			&graph.Description,
			&graph.CreatedAt,
			&graph.UpdatedAt,
			&graph.StartingAt,
			&totalNodes,
			&unlockedNodes,
			&unlockedNonStart,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan graph: %w", err)
		}

		switch {
		case unlockedNonStart == 0:
			graph.Status = models.GraphNotStarted
		case unlockedNodes == totalNodes:
			graph.Status = models.GraphCompleted
		default:
			graph.Status = models.GraphInProgress
		}

		graphs = append(graphs, graph)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating graphs: %w", err)
	}

	return graphs, nil
}

func (s *sqliteDB) DeleteGraph(ctx context.Context, graphID string) error {
	_, err := s.executor().ExecContext(ctx, `
		DELETE FROM graph
		WHERE id = ?
	`, graphID)
	if err != nil {
		return fmt.Errorf("failed to delete graph %s: %w", graphID, err)
	}

	return nil
}

func (s *sqliteDB) CreateGraph(ctx context.Context, req models.NewGraphRequest) (string, error) {
	var graphID int64

	err := s.withTx(ctx, nil, func(txRepo *sqliteDB) error {
		// Insert graph
		result, err := txRepo.executor().ExecContext(ctx, `
			INSERT INTO graph (title, description, starting_at)
			VALUES (?, ?, ?)
		`, req.Title, req.Description, req.StartingAt)
		if err != nil {
			return fmt.Errorf("failed to insert graph: %w", err)
		}

		graphID, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get ID for graph: %w", err)
		}

		// Insert start node
		_, err = txRepo.executor().ExecContext(ctx, `
			INSERT INTO node (id, graph_id, type, title)
			VALUES (?, ?, ?, ?)
		`, crypto.UUIDV4(), graphID, models.StartNode, "Start")
		if err != nil {
			return fmt.Errorf("failed to insert start node: %w", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", graphID), nil
}

func (s *sqliteDB) GetGraph(ctx context.Context, graphID string) (graph *models.Graph, err error) {
	return graph, s.withTx(ctx, &sql.TxOptions{ReadOnly: true}, func(txRepo *sqliteDB) error {
		graph, err = s.getGraph(ctx, graphID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}

		nodes, err := s.getCompleteNodes(ctx, graphID)
		if err != nil {
			return err
		}

		// Get all edges for those nodes
		edges, err := s.getEdges(ctx, graphID)
		if err != nil {
			return err
		}

		graph.Nodes = &nodes
		graph.Edges = &edges
		return nil
	})
}

func (s *sqliteDB) GetAccessibleGraph(ctx context.Context, graphID string) (graph *models.Graph, err error) {
	return graph, s.withTx(ctx, &sql.TxOptions{ReadOnly: true}, func(txRepo *sqliteDB) error {
		graph, err = s.getGraph(ctx, graphID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}

		if graph.StartingAt.After(time.Now()) {
			graph = nil
			return nil
		}

		nodes, err := s.getAccessibleNodes(ctx, graphID)
		if err != nil {
			return err
		}

		edges, err := s.getAccessibleEdges(ctx, graphID)
		if err != nil {
			return err
		}

		graph.Nodes = &nodes
		graph.Edges = &edges
		return nil
	})
}

func (s *sqliteDB) UpdateGraph(ctx context.Context, req models.Graph) error {
	return s.withTx(ctx, nil, func(txRepo *sqliteDB) error {
		// update graph details
		var viewportX, viewportY, viewportZoom sql.NullFloat64
		if req.Viewport != nil {
			viewportX = sql.NullFloat64{Float64: req.Viewport.X, Valid: true}
			viewportY = sql.NullFloat64{Float64: req.Viewport.Y, Valid: true}
			viewportZoom = sql.NullFloat64{Float64: req.Viewport.Zoom, Valid: true}
		}
		result, err := txRepo.executor().ExecContext(ctx, `
			UPDATE graph
			SET title = ?, description = ?, starting_at = ?, viewport_x = ?, viewport_y = ?, viewport_zoom = ?
			WHERE id = ?;
		`, req.Title, req.Description, req.StartingAt, viewportX, viewportY, viewportZoom, req.ID)
		if err != nil {
			return fmt.Errorf("failed to update graph: %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected for graph update: %w", err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("no graph found with ID %s", req.ID)
		}

		err = txRepo.updateNodes(ctx, req.ID, *req.Nodes)
		if err != nil {
			return fmt.Errorf("failed to update nodes: %w", err)
		}

		err = txRepo.updateEdges(ctx, req.ID, *req.Edges)
		if err != nil {
			return fmt.Errorf("failed to update edges: %w", err)
		}

		return nil
	})
}

func (s *sqliteDB) GetAccessibleNode(ctx context.Context, nodeID string) (*models.Node, error) {
	// Look up graph_id from the node itself, then reuse the shared accessible CTEs.
	query := `WITH
       node_graph AS (
          SELECT graph_id FROM node WHERE id = ?1
       ),` + accessibleCTEsFor("(SELECT graph_id FROM node_graph)") + `
       SELECT * FROM node_full
       WHERE id = ?1
         AND id IN (SELECT id FROM accessible);`

	row := s.executor().QueryRowContext(ctx, query, nodeID)

	node, err := s.scanNodeFull(row)
	if err != nil {
		return nil, fmt.Errorf("failed to query accessible node: %w", err)
	}

	node.UIPosition = nil

	return &node, nil
}

func (s *sqliteDB) UnlockNode(ctx context.Context, nodeID string) error {
	_, err := s.executor().ExecContext(ctx, `
		UPDATE node
		SET unlocked_at = CURRENT_TIMESTAMP
		WHERE id = ? AND unlocked_at IS NULL
	`, nodeID)
	if err != nil {
		return fmt.Errorf("failed to unlock node %s: %w", nodeID, err)
	}

	return nil
}

func (s *sqliteDB) LockNode(ctx context.Context, nodeID string) error {
	_, err := s.executor().ExecContext(ctx, `
		UPDATE node
		SET unlocked_at = NULL
		WHERE id = ? AND type != 'start'
	`, nodeID)
	if err != nil {
		return fmt.Errorf("failed to lock node %s: %w", nodeID, err)
	}

	return nil
}
