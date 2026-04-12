package service

import (
	"backend/internal/crypto"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/repository/storage"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

const RegistrationExpiryDuration = 10 * time.Minute

var (
	ErrInvalidGraph = errors.New("invalid graph data")
)

type AdminService struct {
	repo  repository.Repository
	files storage.FileStore
}

func NewAdminService(repo repository.Repository, fileStore storage.FileStore) *AdminService {
	return &AdminService{repo: repo, files: fileStore}
}

func (s *AdminService) CreateRegistrationCode(ctx context.Context) (models.RegistrationCode, error) {
	createdAt := time.Now()
	expiry := createdAt.Add(RegistrationExpiryDuration)

	code, err := crypto.GenerateHumanReadableCode()
	if err != nil {
		return models.RegistrationCode{}, err
	}

	model := models.RegistrationCode{
		Code:      code,
		ExpiresAt: expiry,
		CreatedAt: createdAt,
	}

	if err := s.repo.RefreshRegistrationCode(ctx, model); err != nil {
		return model, err
	}

	return model, nil
}

func (s *AdminService) InvalidateRegistrationCode(ctx context.Context) error {
	return s.repo.DeleteRegistrationCode(ctx)
}

func (s *AdminService) GetRegistrationCode(ctx context.Context) (*models.RegistrationCode, error) {
	return s.repo.GetRegistrationCode(ctx)
}

func (s *AdminService) ListDevices(ctx context.Context) ([]models.Device, error) {
	return s.repo.ListDevices(ctx)
}

func (s *AdminService) RevokeDevice(ctx context.Context, deviceID string) error {
	return s.repo.DeleteDevice(ctx, deviceID)
}

func (s *AdminService) UpdateDeviceInfo(ctx context.Context, deviceID string, deviceInfo string) error {
	return s.repo.UpdateDeviceInfo(ctx, deviceID, deviceInfo)
}

func (s *AdminService) CreateFavourChoice(ctx context.Context, choice *models.FavourChoice) error {
	return s.repo.CreateFavourChoice(ctx, choice)
}

func (s *AdminService) UpdateFavourChoice(ctx context.Context, choice models.FavourChoice) error {
	return s.repo.UpdateFavourChoice(ctx, choice)
}

func (s *AdminService) DeleteFavourChoice(ctx context.Context, choiceID string) error {
	return s.repo.DeleteFavourChoice(ctx, choiceID)
}

func (s *AdminService) UpdateFavourCount(ctx context.Context, count int) error {
	return s.repo.UpdateFavourCount(ctx, count)
}

func (s *AdminService) FulfilFavourRequest(ctx context.Context, requestID string) error {
	return s.repo.FulfilFavourRequest(ctx, requestID)
}
func (s *AdminService) UnfulfilFavourRequest(ctx context.Context, requestID string) error {
	return s.repo.UnfulfilFavourRequest(ctx, requestID)
}

func (s *AdminService) ListGraphs(ctx context.Context) ([]models.Graph, error) {
	return s.repo.ListGraphs(ctx)
}

func (s *AdminService) CreateGraph(ctx context.Context, req models.NewGraphRequest) (string, error) {
	return s.repo.CreateGraph(ctx, req)
}

func (s *AdminService) DeleteGraph(ctx context.Context, graphID string) error {
	return s.repo.DeleteGraph(ctx, graphID)
}

func (s *AdminService) GetGraph(ctx context.Context, graphID string) (*models.Graph, error) {
	graph, err := s.repo.GetGraph(ctx, graphID)
	if err != nil {
		return nil, err
	}

	return graph, nil
}

func (s *AdminService) UpdateGraph(ctx context.Context, graph models.Graph) error {
	nodes := graph.Nodes
	if nodes != nil {
		startNodeSeen := false
		for _, node := range *nodes {
			if node.GraphID != "" && node.GraphID != graph.ID {
				return ErrInvalidGraph
			}

			if node.Type == models.StartNode {
				if startNodeSeen {
					return ErrInvalidGraph
				} else {
					startNodeSeen = true
				}
			}
		}
	}

	edges := graph.Edges
	if edges != nil {
		for _, edge := range *edges {
			if edge.GraphID != "" && edge.GraphID != graph.ID {
				return ErrInvalidGraph
			}
		}
	}

	return s.repo.UpdateGraph(ctx, graph)
}

func (s *AdminService) AdminUnlockNode(ctx context.Context, nodeID string) error {
	return s.repo.UnlockNode(ctx, nodeID)
}

func (s *AdminService) AdminLockNode(ctx context.Context, nodeID string) error {
	return s.repo.LockNode(ctx, nodeID)
}

func (s *AdminService) UploadFile(ctx context.Context, nodeID, filename, mime string, size int64, reader io.ReadSeeker) (*models.File, error) {
	storageKey, err := s.files.Put(ctx, filename, mime, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	file := models.File{
		NodeID:     nodeID,
		Filename:   filename,
		MimeType:   mime,
		SizeBytes:  size,
		StorageKey: storageKey,
		CreatedAt:  time.Now(),
	}

	err = s.repo.CreateFile(ctx, &file)
	if err != nil {
		// Attempt to clean up the uploaded file if database operation fails
		if delErr := s.files.Delete(ctx, storageKey); delErr != nil {
			log.Printf("failed to clean up uploaded file after database error: %v", delErr)
		}
		return nil, fmt.Errorf("failed to create file record in database: %w", err)
	}

	return &file, nil
}

func (s *AdminService) ListFiles(ctx context.Context, nodeID string) ([]models.File, error) {
	return s.repo.ListFiles(ctx, nodeID)
}

//func (s *AdminService) CleanupOrphanedFiles(ctx context.Context) error {
//	allKeys, err := s.files.List(ctx)
//	if err != nil {
//		return fmt.Errorf("failed to list files in store: %w", err)
//	}
//
//	for _, key := range allKeys {
//		referenced, err := s.db.FileExists(ctx, key)
//		if err != nil {
//			return fmt.Errorf("failed to check db for key %s: %w", key, err)
//		}
//		if !referenced {
//			if err := s.fileStore.Delete(ctx, key); err != nil {
//				// log and continue rather than aborting the whole sweep
//				slog.Error("failed to delete orphaned file", "key", key, "err", err)
//			}
//		}
//	}
//	return nil
//}
