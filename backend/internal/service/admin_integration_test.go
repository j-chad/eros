package service

import (
	"backend/internal/models"
	"backend/internal/testutil"
	"backend/internal/testutil/testdb"
	"context"
	"strings"
	"testing"
	"time"
)

func TestUpdateGraph_HappyPath(t *testing.T) {
	repo := testdb.New(t)
	svc := NewAdminService(repo, nil, nil)
	ctx := context.Background()

	graphID, err := repo.CreateGraph(ctx, models.NewGraphRequest{
		Title:      "Test Graph",
		StartingAt: time.Now(),
	})
	testutil.NilErr(t, err)

	graph, err := repo.GetGraph(ctx, graphID)
	testutil.NilErr(t, err)
	startNodeID := (*graph.Nodes)[0].ID

	nodes := []models.Node{
		{ID: startNodeID, Type: models.StartNode, Title: "Start"},
		{ID: "new-node", Type: models.CodeGateNode, Title: "Code Gate"},
	}
	edges := []models.Edge{
		{ID: "e1", From: startNodeID, To: "new-node"},
	}
	updated := models.Graph{
		ID:    graphID,
		Title: "Updated",
		Nodes: &nodes,
		Edges: &edges,
	}

	testutil.NilErr(t, svc.UpdateGraph(ctx, updated))

	got, err := svc.GetGraph(ctx, graphID)
	testutil.NilErr(t, err)
	testutil.Equal(t, got.Title, "Updated")
}

func TestUploadFile_HappyPath(t *testing.T) {
	repo := testdb.New(t)
	store := testdb.NewFileStore(t)
	svc := NewAdminService(repo, store, nil)
	ctx := context.Background()

	graphID, _ := repo.CreateGraph(ctx, models.NewGraphRequest{
		Title:      "Graph",
		StartingAt: time.Now(),
	})
	graph, _ := repo.GetGraph(ctx, graphID)
	startNodeID := (*graph.Nodes)[0].ID

	rewardNodeID := "reward-node-1"
	nodes := []models.Node{
		{ID: startNodeID, Type: models.StartNode, Title: "Start"},
		{ID: rewardNodeID, Type: models.RewardNode, Title: "Reward"},
	}
	edges := []models.Edge{
		{ID: "e1", From: startNodeID, To: rewardNodeID},
	}
	repo.UpdateGraph(ctx, models.Graph{ID: graphID, Nodes: &nodes, Edges: &edges})

	file, err := svc.UploadFile(ctx, rewardNodeID, "photo.jpg", "image/jpeg", 1024, strings.NewReader("image data"))
	testutil.NilErr(t, err)
	testutil.Equal(t, file.Filename, "photo.jpg")
	testutil.Equal(t, file.MimeType, "image/jpeg")
	testutil.Equal(t, file.SizeBytes, int64(1024))
	testutil.True(t, file.ID != "", "file should get an ID")

	// Verify file on disk
	reader, err := store.Get(ctx, file.StorageKey)
	testutil.NilErr(t, err)
	reader.Close()

	// Verify DB record
	files, err := repo.ListFiles(ctx, rewardNodeID)
	testutil.NilErr(t, err)
	testutil.Equal(t, len(files), 1)
}

func TestListGraphs_Sanitization(t *testing.T) {
	repo := testdb.New(t)
	svc := NewGraphService(repo, nil)
	ctx := context.Background()

	repo.CreateGraph(ctx, models.NewGraphRequest{
		Title:       "Past Graph",
		Description: "should be visible",
		StartingAt:  time.Now().Add(-24 * time.Hour),
	})

	repo.CreateGraph(ctx, models.NewGraphRequest{
		Title:       "Future Graph",
		Description: "should be hidden",
		StartingAt:  time.Now().Add(24 * time.Hour),
	})

	graphs, err := svc.ListGraphs(ctx)
	testutil.NilErr(t, err)
	testutil.Equal(t, len(graphs), 2)

	for _, g := range graphs {
		if g.StartingAt.After(time.Now()) {
			testutil.Equal(t, g.Title, "")
			testutil.Equal(t, g.Description, "")
		} else {
			testutil.Equal(t, g.Title, "Past Graph")
			testutil.Equal(t, g.Description, "should be visible")
		}
	}
}

func TestGetGraph_StripsViewport(t *testing.T) {
	repo := testdb.New(t)
	svc := NewGraphService(repo, nil)
	ctx := context.Background()

	graphID, _ := repo.CreateGraph(ctx, models.NewGraphRequest{
		Title:      "Graph",
		StartingAt: time.Now().Add(-time.Hour),
	})

	graph, err := svc.GetGraph(ctx, graphID)
	testutil.NilErr(t, err)
	testutil.Nil(t, graph.Viewport)
}
