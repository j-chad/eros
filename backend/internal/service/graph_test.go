package service

import (
	"backend/internal/models"
	"backend/internal/testutil"
	"testing"
	"time"
)

func TestValidateUnlockPayload_CodeGate_Correct(t *testing.T) {
	node := &models.Node{
		ID:   "n1",
		Type: models.CodeGateNode,
		Data: &models.CodeData{Codes: []string{"SECRET"}},
	}
	testutil.NilErr(t, validateUnlockPayload(node, "SECRET"))
}

func TestValidateUnlockPayload_CodeGate_Wrong(t *testing.T) {
	node := &models.Node{
		ID:   "n1",
		Type: models.CodeGateNode,
		Data: &models.CodeData{Codes: []string{"SECRET"}},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, "WRONG"), ErrNodeUnlockIncorrect)
}

func TestValidateUnlockPayload_CodeGate_Empty(t *testing.T) {
	node := &models.Node{
		ID:   "n1",
		Type: models.CodeGateNode,
		Data: &models.CodeData{Codes: []string{"SECRET"}},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, ""), ErrNodeUnlockIncorrect)
}

func TestValidateUnlockPayload_CodeGate_CaseInsensitive(t *testing.T) {
	node := &models.Node{
		ID:   "n1",
		Type: models.CodeGateNode,
		Data: &models.CodeData{Codes: []string{"Secret"}},
	}
	testutil.NilErr(t, validateUnlockPayload(node, "secret"))
}

func TestValidateUnlockPayload_LocationGate_WithinRadius(t *testing.T) {
	node := &models.Node{
		ID:   "n2",
		Type: models.LocationGateNode,
		Data: &models.LocationData{LocationArea: models.LocationArea{Latitude: -36.8485, Longitude: 174.7633, RadiusM: 100}},
	}
	testutil.NilErr(t, validateUnlockPayload(node, "-36.8484,174.7633"))
}

func TestValidateUnlockPayload_LocationGate_OutsideRadius(t *testing.T) {
	node := &models.Node{
		ID:   "n2",
		Type: models.LocationGateNode,
		Data: &models.LocationData{LocationArea: models.LocationArea{Latitude: -36.8485, Longitude: 174.7633, RadiusM: 10}},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, "-36.8400,174.7633"), ErrNodeUnlockIncorrect)
}

func TestValidateUnlockPayload_LocationGate_DefaultRadius(t *testing.T) {
	node := &models.Node{
		ID:   "n2",
		Type: models.LocationGateNode,
		Data: &models.LocationData{LocationArea: models.LocationArea{Latitude: -36.8485, Longitude: 174.7633, RadiusM: 0}},
	}
	// ~5m away - within default 10m
	testutil.NilErr(t, validateUnlockPayload(node, "-36.84846,174.7633"))
}

func TestValidateUnlockPayload_LocationGate_InvalidPayload(t *testing.T) {
	node := &models.Node{
		ID:   "n2",
		Type: models.LocationGateNode,
		Data: &models.LocationData{LocationArea: models.LocationArea{Latitude: 0, Longitude: 0, RadiusM: 10}},
	}

	tests := []struct {
		name    string
		payload string
	}{
		{"empty", ""},
		{"no comma", "123.456"},
		{"letters", "abc,def"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUnlockPayload(node, tt.payload)
			testutil.NotNilErr(t, err)
			testutil.ErrorIsNot(t, err, ErrNodeUnlockIncorrect)
		})
	}
}

func TestValidateUnlockPayload_ManualGate_Unlocked(t *testing.T) {
	now := time.Now()
	node := &models.Node{
		ID:   "n3",
		Type: models.ManualNode,
		Data: &models.ManualData{Instructions: "do the thing", UnlockedAt: &now},
	}
	testutil.NilErr(t, validateUnlockPayload(node, ""))
}

func TestValidateUnlockPayload_ManualGate_NotUnlocked(t *testing.T) {
	node := &models.Node{
		ID:   "n3",
		Type: models.ManualNode,
		Data: &models.ManualData{Instructions: "do the thing", UnlockedAt: nil},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, ""), ErrNodeUnlockIncorrect)
}

func TestValidateUnlockPayload_TimeGate_Past(t *testing.T) {
	node := &models.Node{
		ID:   "n7",
		Type: models.TimeGateNode,
		Data: &models.TimeData{UnlockAt: time.Now().Add(-1 * time.Hour)},
	}
	testutil.NilErr(t, validateUnlockPayload(node, ""))
}

func TestValidateUnlockPayload_TimeGate_Future(t *testing.T) {
	node := &models.Node{
		ID:   "n7",
		Type: models.TimeGateNode,
		Data: &models.TimeData{UnlockAt: time.Now().Add(1 * time.Hour)},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, ""), ErrNodeUnlockIncorrect)
}

func TestValidateUnlockPayload_UnsupportedType(t *testing.T) {
	node := &models.Node{ID: "n4", Type: models.StartNode}
	testutil.NotNilErr(t, validateUnlockPayload(node, ""))
}

func TestValidateUnlockPayload_RewardType(t *testing.T) {
	node := &models.Node{ID: "n5", Type: models.RewardNode}
	testutil.NotNilErr(t, validateUnlockPayload(node, ""))
}

func TestValidateUnlockPayload_WrongDataType(t *testing.T) {
	node := &models.Node{
		ID:   "n6",
		Type: models.CodeGateNode,
		Data: &models.LocationData{LocationArea: models.LocationArea{Latitude: 0, Longitude: 0}},
	}
	testutil.NotNilErr(t, validateUnlockPayload(node, "anything"))
}

func TestUpdateGraph_DuplicateStartNodes(t *testing.T) {
	nodes := []models.Node{
		{ID: "1", Type: models.StartNode},
		{ID: "2", Type: models.StartNode},
	}
	graph := models.Graph{ID: "g1", Nodes: &nodes}

	svc := &GraphService{repo: nil}
	testutil.ErrorIs(t, svc.UpdateGraph(nil, graph), ErrInvalidGraph)
}

func TestUpdateGraph_NodeGraphIDMismatch(t *testing.T) {
	nodes := []models.Node{
		{ID: "1", Type: models.StartNode, GraphID: "other-graph"},
	}
	graph := models.Graph{ID: "g1", Nodes: &nodes}

	svc := &GraphService{repo: nil}
	testutil.ErrorIs(t, svc.UpdateGraph(nil, graph), ErrInvalidGraph)
}

func TestUpdateGraph_EdgeGraphIDMismatch(t *testing.T) {
	edges := []models.Edge{
		{ID: "e1", GraphID: "other-graph", From: "a", To: "b"},
	}
	graph := models.Graph{ID: "g1", Edges: &edges}

	svc := &GraphService{repo: nil}
	testutil.ErrorIs(t, svc.UpdateGraph(nil, graph), ErrInvalidGraph)
}

func TestUpdateGraph_NodeEmptyGraphID_PassesValidation(t *testing.T) {
	nodes := []models.Node{
		{ID: "1", Type: models.StartNode, GraphID: ""},
		{ID: "2", Type: models.CodeGateNode, GraphID: ""},
	}
	graph := models.Graph{ID: "g1", Nodes: &nodes}

	svc := &GraphService{repo: nil}
	// Validation passes, then panics on nil repo - recover to confirm validation OK
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

	svc := &GraphService{repo: nil}
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

	svc := &GraphService{repo: nil}
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

	svc := &GraphService{repo: nil}
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
