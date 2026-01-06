package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *sqliteDB) ListGraphs(ctx context.Context) ([]models.Graph, error) {
	query := `
        SELECT 
            g.id,
            g.title,
            g.description,
            g.created_at,
            g.updated_at,
            g.starting_at
        FROM graph g
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

		err := rows.Scan(
			&graph.ID,
			&graph.Title,
			&graph.Description,
			&graph.CreatedAt,
			&graph.UpdatedAt,
			&graph.StartingAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan graph: %w", err)
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

	err := s.withTx(ctx, func(txRepo *sqliteDB) error {
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
			INSERT INTO node (graph_id, type, title)
			VALUES (?, ?, ?)
		`, graphID, models.StartNode, "Start")
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

func (s *sqliteDB) GetGraph(ctx context.Context, graphID string) (*models.Graph, error) {
	row := s.executor().QueryRowContext(ctx, `
		SELECT 
		    title,
		    description,
		    starting_at,
		    viewport_x,
		    viewport_y,
		    viewport_zoom,
		    created_at,
		    updated_at
		FROM graph
		WHERE id = ?;
	`, graphID, models.StartNode)

	graph := &models.Graph{ID: graphID}
	var viewportX, viewportY, viewportZoom sql.NullFloat64

	err := row.Scan(&graph.Title, &graph.Description, &graph.StartingAt, &viewportX, &viewportY, &viewportZoom, &graph.CreatedAt, &graph.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan graph: %w", err)
	}

	if viewportX.Valid && viewportY.Valid && viewportZoom.Valid {
		graph.Viewport = &models.Viewport{
			X:    viewportX.Float64,
			Y:    viewportY.Float64,
			Zoom: viewportZoom.Float64,
		}
	}

	nodes, err := s.getCompleteNodes(ctx, graphID)
	if err != nil {
		return nil, err
	}

	// Get all edges for those nodes
	edges, err := s.getEdges(ctx, graphID)
	if err != nil {
		return nil, err
	}

	graph.Nodes = &nodes
	graph.Edges = &edges

	return graph, nil
}

func (s *sqliteDB) UpdateGraph(ctx context.Context, req models.Graph) error {
	return s.withTx(ctx, func(txRepo *sqliteDB) error {
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

		err = txRepo.updateEdges(ctx, req.ID, *req.Edges)
		if err != nil {
			return fmt.Errorf("failed to update edges: %w", err)
		}

		err = txRepo.updateNodes(ctx, req.ID, *req.Nodes)
		if err != nil {
			return fmt.Errorf("failed to update nodes: %w", err)
		}

		return nil
	})
}
