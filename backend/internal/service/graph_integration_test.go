//go:build integration

package service

import (
	"backend/internal/models"
	"backend/internal/testutil"
	"backend/internal/testutil/testdb"
	"context"
	"testing"
	"time"
)

// buildGraph is a helper that creates a graph with the given nodes and edges.
// Use "start" as the start node ID in edges; it is replaced with the real ID.
func buildGraph(t *testing.T, repo interface {
	CreateGraph(ctx context.Context, req models.NewGraphRequest) (string, error)
	GetGraph(ctx context.Context, graphID string) (*models.Graph, error)
	UpdateGraph(ctx context.Context, graph models.Graph) error
}, nodes []models.Node, edges []models.Edge) string {
	t.Helper()
	ctx := context.Background()

	graphID, err := repo.CreateGraph(ctx, models.NewGraphRequest{
		Title:      "Test Graph",
		StartingAt: time.Now().Add(-time.Hour),
	})
	testutil.NilErr(t, err)

	graph, err := repo.GetGraph(ctx, graphID)
	testutil.NilErr(t, err)
	startNodeID := (*graph.Nodes)[0].ID

	// Prepend the start node to the provided nodes.
	allNodes := append([]models.Node{{ID: startNodeID, Type: models.StartNode, Title: "Start"}}, nodes...)

	// Replace "start" placeholder in edges with the real start node ID.
	for i := range edges {
		if edges[i].From == "start" {
			edges[i].From = startNodeID
		}
		if edges[i].To == "start" {
			edges[i].To = startNodeID
		}
	}

	testutil.NilErr(t, repo.UpdateGraph(ctx, models.Graph{
		ID:    graphID,
		Nodes: &allNodes,
		Edges: &edges,
	}))

	return graphID
}

// findNode returns the first node with the given ID from a slice, or nil.
func findNode(nodes []models.Node, id string) *models.Node {
	for i := range nodes {
		if nodes[i].ID == id {
			return &nodes[i]
		}
	}
	return nil
}

// nodeIDs returns a sorted set of node IDs from a slice.
func nodeIDs(nodes []models.Node) []string {
	ids := make([]string, len(nodes))
	for i, n := range nodes {
		ids[i] = n.ID
	}
	return ids
}

func TestUnlockNode_AutoUnlocksRewardNode(t *testing.T) {
	repo := testdb.New(t)
	svc := NewGraphService(repo, NewFileService(repo, testdb.NewFileStore(t)))
	ctx := context.Background()

	// Graph: Start → LocationGate → Reward
	buildGraph(t, repo,
		[]models.Node{
			{ID: "gate1", Type: models.LocationGateNode, Title: "Gate", Data: models.LocationData{LocationArea: models.LocationArea{Latitude: 0, Longitude: 0, RadiusM: 1000000}}},
			{ID: "reward1", Type: models.RewardNode, Title: "Reward", Data: models.RewardData{RewardType: "image", Payload: "photo.jpg"}},
		},
		[]models.Edge{
			{ID: "e1", From: "start", To: "gate1"},
			{ID: "e2", From: "gate1", To: "reward1"},
		},
	)

	result, err := svc.UnlockNode(ctx, "gate1", "0,0")
	testutil.NilErr(t, err)

	// The reward node should appear in new_nodes with unlocked_at set.
	reward := findNode(result.NewNodes, "reward1")
	testutil.NotNil(t, reward)
	testutil.NotNil(t, reward.UnlockedAt)
}

func TestUnlockNode_AutoUnlocksChainedRewards(t *testing.T) {
	repo := testdb.New(t)
	svc := NewGraphService(repo, NewFileService(repo, testdb.NewFileStore(t)))
	ctx := context.Background()

	// Graph: Start → LocationGate → Reward1 → Reward2 → Reward3
	buildGraph(t, repo,
		[]models.Node{
			{ID: "gate1", Type: models.LocationGateNode, Title: "Gate", Data: models.LocationData{LocationArea: models.LocationArea{Latitude: 0, Longitude: 0, RadiusM: 1000000}}},
			{ID: "reward1", Type: models.RewardNode, Title: "Reward 1", Data: models.RewardData{RewardType: "image", Payload: "a.jpg"}},
			{ID: "reward2", Type: models.RewardNode, Title: "Reward 2", Data: models.RewardData{RewardType: "markdown", Payload: "hello"}},
			{ID: "reward3", Type: models.RewardNode, Title: "Reward 3", Data: models.RewardData{RewardType: "favour", Payload: ""}},
		},
		[]models.Edge{
			{ID: "e1", From: "start", To: "gate1"},
			{ID: "e2", From: "gate1", To: "reward1"},
			{ID: "e3", From: "reward1", To: "reward2"},
			{ID: "e4", From: "reward2", To: "reward3"},
		},
	)

	result, err := svc.UnlockNode(ctx, "gate1", "0,0")
	testutil.NilErr(t, err)

	// All three chained rewards should be present and auto-unlocked.
	for _, id := range []string{"reward1", "reward2", "reward3"} {
		n := findNode(result.NewNodes, id)
		testutil.NotNil(t, n)
		testutil.NotNil(t, n.UnlockedAt)
	}
}

func TestUnlockNode_ChainedRewardThenGate(t *testing.T) {
	repo := testdb.New(t)
	svc := NewGraphService(repo, NewFileService(repo, testdb.NewFileStore(t)))
	ctx := context.Background()

	// Graph: Start → LocationGate1 → Reward → LocationGate2
	// The reward should auto-unlock, but LocationGate2 should NOT.
	buildGraph(t, repo,
		[]models.Node{
			{ID: "gate1", Type: models.LocationGateNode, Title: "Gate 1", Data: models.LocationData{LocationArea: models.LocationArea{Latitude: 0, Longitude: 0, RadiusM: 1000000}}},
			{ID: "reward1", Type: models.RewardNode, Title: "Reward", Data: models.RewardData{RewardType: "image", Payload: "a.jpg"}},
			{ID: "gate2", Type: models.LocationGateNode, Title: "Gate 2", Data: models.LocationData{LocationArea: models.LocationArea{Latitude: 10, Longitude: 10, RadiusM: 100}}},
		},
		[]models.Edge{
			{ID: "e1", From: "start", To: "gate1"},
			{ID: "e2", From: "gate1", To: "reward1"},
			{ID: "e3", From: "reward1", To: "gate2"},
		},
	)

	result, err := svc.UnlockNode(ctx, "gate1", "0,0")
	testutil.NilErr(t, err)

	// Reward should be auto-unlocked.
	reward := findNode(result.NewNodes, "reward1")
	testutil.NotNil(t, reward)
	testutil.NotNil(t, reward.UnlockedAt)

	// Gate2 should be accessible but NOT unlocked.
	gate2 := findNode(result.NewNodes, "gate2")
	testutil.NotNil(t, gate2)
	testutil.Nil(t, gate2.UnlockedAt)
}

func TestUnlockNode_NoReward_NoAutoUnlock(t *testing.T) {
	repo := testdb.New(t)
	svc := NewGraphService(repo, NewFileService(repo, testdb.NewFileStore(t)))
	ctx := context.Background()

	// Graph: Start → LocationGate1 → LocationGate2
	// No rewards - gate2 should be accessible but not unlocked.
	buildGraph(t, repo,
		[]models.Node{
			{ID: "gate1", Type: models.LocationGateNode, Title: "Gate 1", Data: models.LocationData{LocationArea: models.LocationArea{Latitude: 0, Longitude: 0, RadiusM: 1000000}}},
			{ID: "gate2", Type: models.LocationGateNode, Title: "Gate 2", Data: models.LocationData{LocationArea: models.LocationArea{Latitude: 10, Longitude: 10, RadiusM: 100}}},
		},
		[]models.Edge{
			{ID: "e1", From: "start", To: "gate1"},
			{ID: "e2", From: "gate1", To: "gate2"},
		},
	)

	result, err := svc.UnlockNode(ctx, "gate1", "0,0")
	testutil.NilErr(t, err)

	gate2 := findNode(result.NewNodes, "gate2")
	testutil.NotNil(t, gate2)
	testutil.Nil(t, gate2.UnlockedAt)
}

func TestUnlockNode_GrantsFavourPoints(t *testing.T) {
	repo := testdb.New(t)
	svc := NewGraphService(repo, NewFileService(repo, testdb.NewFileStore(t)))
	ctx := context.Background()

	// Graph: Start → Gate → Reward (give_favours: 3)
	buildGraph(t, repo,
		[]models.Node{
			{ID: "gate1", Type: models.LocationGateNode, Title: "Gate", Data: models.LocationData{LocationArea: models.LocationArea{Latitude: 0, Longitude: 0, RadiusM: 1000000}}},
			{ID: "reward1", Type: models.RewardNode, Title: "Reward", Data: models.RewardData{RewardType: "favour", GiveFavours: 3}},
		},
		[]models.Edge{
			{ID: "e1", From: "start", To: "gate1"},
			{ID: "e2", From: "gate1", To: "reward1"},
		},
	)

	// Verify starting balance is 0.
	countBefore, err := repo.GetFavourCount(ctx)
	testutil.NilErr(t, err)
	testutil.Equal(t, countBefore.Total, 0)

	_, err = svc.UnlockNode(ctx, "gate1", "0,0")
	testutil.NilErr(t, err)

	// Balance should now be 3.
	countAfter, err := repo.GetFavourCount(ctx)
	testutil.NilErr(t, err)
	testutil.Equal(t, countAfter.Total, 3)
	testutil.Equal(t, countAfter.Remaining, 3)
}

func TestUnlockNode_GrantsFavoursFromChainedRewards(t *testing.T) {
	repo := testdb.New(t)
	svc := NewGraphService(repo, NewFileService(repo, testdb.NewFileStore(t)))
	ctx := context.Background()

	// Graph: Start → Gate → Reward1 (2pts) → Reward2 (5pts)
	buildGraph(t, repo,
		[]models.Node{
			{ID: "gate1", Type: models.LocationGateNode, Title: "Gate", Data: models.LocationData{LocationArea: models.LocationArea{Latitude: 0, Longitude: 0, RadiusM: 1000000}}},
			{ID: "reward1", Type: models.RewardNode, Title: "Reward 1", Data: models.RewardData{RewardType: "image", Payload: "a.jpg", GiveFavours: 2}},
			{ID: "reward2", Type: models.RewardNode, Title: "Reward 2", Data: models.RewardData{RewardType: "favour", GiveFavours: 5}},
		},
		[]models.Edge{
			{ID: "e1", From: "start", To: "gate1"},
			{ID: "e2", From: "gate1", To: "reward1"},
			{ID: "e3", From: "reward1", To: "reward2"},
		},
	)

	_, err := svc.UnlockNode(ctx, "gate1", "0,0")
	testutil.NilErr(t, err)

	// Should be 2 + 5 = 7.
	count, err := repo.GetFavourCount(ctx)
	testutil.NilErr(t, err)
	testutil.Equal(t, count.Total, 7)
	testutil.Equal(t, count.Remaining, 7)
}

func TestUnlockNode_NoFavoursWhenGiveFavoursZero(t *testing.T) {
	repo := testdb.New(t)
	svc := NewGraphService(repo, NewFileService(repo, testdb.NewFileStore(t)))
	ctx := context.Background()

	// Graph: Start → Gate → Reward (give_favours: 0)
	buildGraph(t, repo,
		[]models.Node{
			{ID: "gate1", Type: models.LocationGateNode, Title: "Gate", Data: models.LocationData{LocationArea: models.LocationArea{Latitude: 0, Longitude: 0, RadiusM: 1000000}}},
			{ID: "reward1", Type: models.RewardNode, Title: "Reward", Data: models.RewardData{RewardType: "image", Payload: "a.jpg", GiveFavours: 0}},
		},
		[]models.Edge{
			{ID: "e1", From: "start", To: "gate1"},
			{ID: "e2", From: "gate1", To: "reward1"},
		},
	)

	_, err := svc.UnlockNode(ctx, "gate1", "0,0")
	testutil.NilErr(t, err)

	// Balance should remain 0.
	count, err := repo.GetFavourCount(ctx)
	testutil.NilErr(t, err)
	testutil.Equal(t, count.Total, 0)
}
