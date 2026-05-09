//go:build integration

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

func TestUploadFile_ReplacesExisting(t *testing.T) {
	repo := testdb.New(t)
	store := testdb.NewFileStore(t)
	svc := NewAdminService(repo, store, NewFileService(repo, store))
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
		{ID: rewardNodeID, Type: models.RewardNode, Title: "Reward", Data: models.RewardData{RewardType: "image"}},
	}
	edges := []models.Edge{{ID: "e1", From: startNodeID, To: rewardNodeID}}
	repo.UpdateGraph(ctx, models.Graph{ID: graphID, Nodes: &nodes, Edges: &edges})

	// Upload first file.
	file1, err := svc.UploadFile(ctx, rewardNodeID, "first.jpg", "image/jpeg", 100, strings.NewReader("first"))
	testutil.NilErr(t, err)
	testutil.Equal(t, file1.Filename, "first.jpg")

	// Upload second file - should replace the first.
	file2, err := svc.UploadFile(ctx, rewardNodeID, "second.jpg", "image/jpeg", 200, strings.NewReader("second"))
	testutil.NilErr(t, err)
	testutil.Equal(t, file2.Filename, "second.jpg")

	// Only one file should exist for the node.
	files, err := repo.ListFiles(ctx, rewardNodeID)
	testutil.NilErr(t, err)
	testutil.Equal(t, len(files), 1)
	testutil.Equal(t, files[0].Filename, "second.jpg")

	// Old storage file should be gone.
	_, err = store.Get(ctx, file1.StorageKey)
	testutil.NotNilErr(t, err)
}

func TestGetFile_ByID(t *testing.T) {
	repo := testdb.New(t)
	store := testdb.NewFileStore(t)
	svc := NewAdminService(repo, store, NewFileService(repo, store))
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
		{ID: rewardNodeID, Type: models.RewardNode, Title: "Reward", Data: models.RewardData{RewardType: "image"}},
	}
	edges := []models.Edge{{ID: "e1", From: startNodeID, To: rewardNodeID}}
	repo.UpdateGraph(ctx, models.Graph{ID: graphID, Nodes: &nodes, Edges: &edges})

	file, err := svc.UploadFile(ctx, rewardNodeID, "photo.jpg", "image/jpeg", 1024, strings.NewReader("data"))
	testutil.NilErr(t, err)

	// GetFile by ID
	got, err := repo.GetFile(ctx, file.ID)
	testutil.NilErr(t, err)
	testutil.NotNil(t, got)
	testutil.Equal(t, got.Filename, "photo.jpg")

	// GetFileByNodeID
	got2, err := repo.GetFileByNodeID(ctx, rewardNodeID)
	testutil.NilErr(t, err)
	testutil.NotNil(t, got2)
	testutil.Equal(t, got2.ID, file.ID)

	// GetFile for nonexistent ID returns nil
	got3, err := repo.GetFile(ctx, "99999")
	testutil.NilErr(t, err)
	testutil.Nil(t, got3)
}

func TestGetGraph_IncludesFileMetadata(t *testing.T) {
	repo := testdb.New(t)
	store := testdb.NewFileStore(t)
	adminSvc := NewAdminService(repo, store, NewFileService(repo, store))
	fileSvc := NewFileService(repo, store)
	graphSvc := NewGraphService(repo, fileSvc)
	ctx := context.Background()

	graphID, _ := repo.CreateGraph(ctx, models.NewGraphRequest{
		Title:      "Graph",
		StartingAt: time.Now().Add(-time.Hour),
	})
	graph, _ := repo.GetGraph(ctx, graphID)
	startNodeID := (*graph.Nodes)[0].ID

	rewardNodeID := "reward-node-1"
	nodes := []models.Node{
		{ID: startNodeID, Type: models.StartNode, Title: "Start"},
		{ID: rewardNodeID, Type: models.RewardNode, Title: "Reward", Data: models.RewardData{RewardType: "image"}},
	}
	edges := []models.Edge{{ID: "e1", From: startNodeID, To: rewardNodeID}}
	repo.UpdateGraph(ctx, models.Graph{ID: graphID, Nodes: &nodes, Edges: &edges})

	// Unlock the reward node so it's accessible.
	repo.UnlockNode(ctx, rewardNodeID)

	// Upload a file to the reward node.
	_, err := adminSvc.UploadFile(ctx, rewardNodeID, "sunset.jpg", "image/jpeg", 2048, strings.NewReader("image"))
	testutil.NilErr(t, err)

	// Get graph via the client-facing service.
	got, err := graphSvc.GetGraph(ctx, graphID)
	testutil.NilErr(t, err)

	// Find the reward node and check file metadata.
	var rewardNode *models.Node
	for _, n := range *got.Nodes {
		if n.ID == rewardNodeID {
			rewardNode = &n
			break
		}
	}
	testutil.NotNil(t, rewardNode)

	rd := testutil.TypeAssert[*models.RewardData](t, rewardNode.Data)
	testutil.NotNil(t, rd.File)
	testutil.Equal(t, rd.File.Filename, "sunset.jpg")
	testutil.Equal(t, rd.File.MimeType, "image/jpeg")
	testutil.Equal(t, rd.File.SizeBytes, int64(2048))
	testutil.True(t, strings.HasPrefix(rd.File.URL, "/api/files/"), "URL should be a local path")
	testutil.Nil(t, rd.File.URLExpires)
}

func TestDeleteFilesByNodeID(t *testing.T) {
	repo := testdb.New(t)
	store := testdb.NewFileStore(t)
	svc := NewAdminService(repo, store, NewFileService(repo, store))
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
		{ID: rewardNodeID, Type: models.RewardNode, Title: "Reward", Data: models.RewardData{RewardType: "file"}},
	}
	edges := []models.Edge{{ID: "e1", From: startNodeID, To: rewardNodeID}}
	repo.UpdateGraph(ctx, models.Graph{ID: graphID, Nodes: &nodes, Edges: &edges})

	file, err := svc.UploadFile(ctx, rewardNodeID, "doc.pdf", "application/pdf", 500, strings.NewReader("pdf"))
	testutil.NilErr(t, err)

	keys, err := repo.DeleteFilesByNodeID(ctx, rewardNodeID)
	testutil.NilErr(t, err)
	testutil.Equal(t, len(keys), 1)
	testutil.Equal(t, keys[0], file.StorageKey)

	// DB should have no files.
	files, err := repo.ListFiles(ctx, rewardNodeID)
	testutil.NilErr(t, err)
	testutil.Equal(t, len(files), 0)
}

func TestGetFilesByNodeIDs(t *testing.T) {
	repo := testdb.New(t)
	store := testdb.NewFileStore(t)
	svc := NewAdminService(repo, store, NewFileService(repo, store))
	ctx := context.Background()

	graphID, _ := repo.CreateGraph(ctx, models.NewGraphRequest{
		Title:      "Graph",
		StartingAt: time.Now(),
	})
	graph, _ := repo.GetGraph(ctx, graphID)
	startNodeID := (*graph.Nodes)[0].ID

	nodes := []models.Node{
		{ID: startNodeID, Type: models.StartNode, Title: "Start"},
		{ID: "gate1", Type: models.CodeGateNode, Title: "Gate", Data: models.CodeData{Codes: []string{"x"}}},
		{ID: "r1", Type: models.RewardNode, Title: "R1", Data: models.RewardData{RewardType: "image"}},
		{ID: "r2", Type: models.RewardNode, Title: "R2", Data: models.RewardData{RewardType: "video"}},
	}
	edges := []models.Edge{
		{ID: "e1", From: startNodeID, To: "gate1"},
		{ID: "e2", From: "gate1", To: "r1"},
		{ID: "e3", From: "r1", To: "r2"},
	}
	err := repo.UpdateGraph(ctx, models.Graph{ID: graphID, Nodes: &nodes, Edges: &edges})
	testutil.NilErr(t, err)

	_, err = svc.UploadFile(ctx, "r1", "a.jpg", "image/jpeg", 100, strings.NewReader("aaaa"))
	testutil.NilErr(t, err)
	_, err = svc.UploadFile(ctx, "r2", "b.mp4", "video/mp4", 200, strings.NewReader("bb"))
	testutil.NilErr(t, err)

	fileMap, err := repo.GetFilesByNodeIDs(ctx, []string{"r1", "r2", "nonexistent"})
	testutil.NilErr(t, err)
	testutil.Equal(t, len(fileMap), 2)
	testutil.Equal(t, fileMap["r1"].Filename, "a.jpg")
	testutil.Equal(t, fileMap["r2"].Filename, "b.mp4")
}
