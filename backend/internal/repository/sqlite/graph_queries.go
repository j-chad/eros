package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// accessibleCTEsFor returns two CTE definitions that determine which nodes are
// reachable by the client. The result starts with a newline and is intended to
// follow a WITH keyword - e.g. `WITH` + accessibleCTEsFor("?1") or after a
// preceding CTE with a trailing comma.
//
//  1. unlocked - every node that is already unlocked or is a start node
//     (start nodes are always considered unlocked).
//
//  2. accessible - the unlocked set, plus destinations that are exactly one
//     hop away from an unlocked source node. A destination is excluded when a
//     *sibling* edge from the same source already has its destination unlocked,
//     meaning the user committed to that other branch and cannot backtrack.
//
// graphIDExpr is a SQL expression that resolves to the graph ID - typically
// "?1" for direct queries or a subquery like "(SELECT graph_id FROM node_graph)".
func accessibleCTEsFor(graphIDExpr string) string {
	return strings.ReplaceAll(`
	unlocked AS (
		SELECT id
		FROM node
		WHERE graph_id = {graph_id}
		  AND (unlocked_at IS NOT NULL OR type = 'start')
	),
	accessible AS (
		SELECT id FROM unlocked
		UNION
		SELECT e.destination_node_id
		FROM edge e
			JOIN unlocked u ON u.id = e.source_node_id
		WHERE e.graph_id = {graph_id}
		  -- Prevent backtracking: if any sibling edge's destination is already
		  -- unlocked the user has committed to that branch, so exclude this one.
		  AND NOT EXISTS (
			SELECT 1 FROM edge sibling
				JOIN node n ON n.id = sibling.destination_node_id
			WHERE sibling.source_node_id = e.source_node_id
			  AND sibling.id != e.id
			  AND n.unlocked_at IS NOT NULL
		  )
	) `, "{graph_id}", graphIDExpr)
}

// getGraph retrieves a graph. It does not populate the nodes or edges.
func (s *sqliteDB) getGraph(ctx context.Context, graphID string) (*models.Graph, error) {
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
	`, graphID)

	graph := &models.Graph{ID: graphID}
	var viewportX, viewportY, viewportZoom sql.NullFloat64

	err := row.Scan(&graph.Title, &graph.Description, &graph.StartingAt, &viewportX, &viewportY, &viewportZoom, &graph.CreatedAt, &graph.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan graph: %w", err)
	}

	if viewportX.Valid && viewportY.Valid && viewportZoom.Valid {
		graph.Viewport = &models.Viewport{
			X:    viewportX.Float64,
			Y:    viewportY.Float64,
			Zoom: viewportZoom.Float64,
		}
	}

	return graph, nil
}

// scanNodeFull scans a row from node_full view into a Node struct
func (s *sqliteDB) scanNodeFull(scanner interface {
	Scan(dest ...any) error
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
	case models.ManualNode:
		node.Data = &models.ManualData{}
		if err := json.Unmarshal(dataJSON, &node.Data); err != nil {
			return node, fmt.Errorf("failed to unmarshal manual data: %w", err)
		}
	case models.TimeGateNode:
		node.Data = &models.TimeData{}
		if err := json.Unmarshal(dataJSON, &node.Data); err != nil {
			return node, fmt.Errorf("failed to unmarshal time data: %w", err)
		}
	case models.StartNode:
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

// getAccessibleNodes retrieves all nodes reachable from a start node that are unlocked or have an edge from an unlocked node
func (s *sqliteDB) getAccessibleNodes(ctx context.Context, graphID string) ([]models.Node, error) {
	query := `WITH` + accessibleCTEsFor("?1") + `
		SELECT * FROM node_full
		WHERE graph_id = ?1
		  AND id IN (SELECT id FROM accessible);`

	rows, err := s.executor().QueryContext(ctx, query, graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to query accessible nodes: %w", err)
	}
	defer rows.Close()

	nodes := make([]models.Node, 0)

	for rows.Next() {
		node, err := s.scanNodeFull(rows)
		if err != nil {
			return nil, err
		}

		node.UIPosition = nil

		nodes = append(nodes, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating accessible nodes: %w", err)
	}

	return nodes, nil
}

// getAccessibleEdges retrieves all edges that have a source node that is unlocked or has an edge from an unlocked node
func (s *sqliteDB) getAccessibleEdges(ctx context.Context, graphID string) ([]models.Edge, error) {
	query := `WITH` + accessibleCTEsFor("?1") + `
	SELECT DISTINCT e.id,
		e.source_node_id,
		e.destination_node_id,
		e.choice_label,
		e.created_at,
		e.updated_at
	FROM edge e
		JOIN accessible a ON a.id = e.destination_node_id
	WHERE e.graph_id = ?1;`

	rows, err := s.executor().QueryContext(ctx, query, graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to query accessible edges: %w", err)
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
		return nil, fmt.Errorf("error iterating accessible edges: %w", err)
	}

	return edges, nil
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
