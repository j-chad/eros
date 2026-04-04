package service

import (
	"backend/internal/models"
	"backend/internal/testutil"
	"testing"
)

func TestUpdateGraph_DuplicateStartNodes(t *testing.T) {
	nodes := []models.Node{
		{ID: "1", Type: models.StartNode},
		{ID: "2", Type: models.StartNode},
	}
	graph := models.Graph{ID: "g1", Nodes: &nodes}

	svc := &AdminService{repo: nil}
	testutil.ErrorIs(t, svc.UpdateGraph(nil, graph), ErrInvalidGraph)
}

func TestUpdateGraph_NodeGraphIDMismatch(t *testing.T) {
	nodes := []models.Node{
		{ID: "1", Type: models.StartNode, GraphID: "other-graph"},
	}
	graph := models.Graph{ID: "g1", Nodes: &nodes}

	svc := &AdminService{repo: nil}
	testutil.ErrorIs(t, svc.UpdateGraph(nil, graph), ErrInvalidGraph)
}

func TestUpdateGraph_EdgeGraphIDMismatch(t *testing.T) {
	edges := []models.Edge{
		{ID: "e1", GraphID: "other-graph", From: "a", To: "b"},
	}
	graph := models.Graph{ID: "g1", Edges: &edges}

	svc := &AdminService{repo: nil}
	testutil.ErrorIs(t, svc.UpdateGraph(nil, graph), ErrInvalidGraph)
}

func TestUpdateGraph_NodeEmptyGraphID_PassesValidation(t *testing.T) {
	nodes := []models.Node{
		{ID: "1", Type: models.StartNode, GraphID: ""},
		{ID: "2", Type: models.CodeGateNode, GraphID: ""},
	}
	graph := models.Graph{ID: "g1", Nodes: &nodes}

	svc := &AdminService{repo: nil}
	// Validation passes, then panics on nil repo — recover to confirm validation OK
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		svc.UpdateGraph(nil, graph)
	}()
	testutil.True(t, panicked, "should panic on nil repo after passing validation")
}

func TestUpdateGraph_NilNodesAndEdges_PassesValidation(t *testing.T) {
	graph := models.Graph{ID: "g1", Nodes: nil, Edges: nil}

	svc := &AdminService{repo: nil}
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		svc.UpdateGraph(nil, graph)
	}()
	testutil.True(t, panicked, "should panic on nil repo after passing validation")
}

func TestUpdateGraph_EdgeEmptyGraphID_PassesValidation(t *testing.T) {
	edges := []models.Edge{
		{ID: "e1", GraphID: "", From: "a", To: "b"},
	}
	graph := models.Graph{ID: "g1", Edges: &edges}

	svc := &AdminService{repo: nil}
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		svc.UpdateGraph(nil, graph)
	}()
	testutil.True(t, panicked, "should panic on nil repo after passing validation")
}

func TestUpdateGraph_MultipleNodesOneStart(t *testing.T) {
	nodes := []models.Node{
		{ID: "1", Type: models.StartNode},
		{ID: "2", Type: models.CodeGateNode},
		{ID: "3", Type: models.LocationGateNode},
		{ID: "4", Type: models.RewardNode},
	}
	graph := models.Graph{ID: "g1", Nodes: &nodes}

	svc := &AdminService{repo: nil}
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		svc.UpdateGraph(nil, graph)
	}()
	testutil.True(t, panicked, "should pass validation with exactly one start node")
}
