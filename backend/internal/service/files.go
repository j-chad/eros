package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/repository/storage"
	"context"
	"fmt"
	"io"
	"time"
)

const presignTTL = 15 * time.Minute

type FileService struct {
	repo  repository.Repository
	files storage.FileStore
}

func NewFileService(repo repository.Repository, files storage.FileStore) *FileService {
	return &FileService{repo: repo, files: files}
}

// GetFile fetches file metadata from the database.
func (s *FileService) GetFile(ctx context.Context, fileID string) (*models.File, error) {
	return s.repo.GetFile(ctx, fileID)
}

// GetFileStream returns the file bytes and metadata. It does not check access - the caller must do that.
func (s *FileService) GetFileStream(ctx context.Context, fileID string) (io.ReadCloser, *models.File, error) {
	file, err := s.repo.GetFile(ctx, fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file: %w", err)
	}
	if file == nil {
		return nil, nil, nil
	}

	reader, err := s.files.Get(ctx, file.StorageKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file from storage: %w", err)
	}

	return reader, file, nil
}

// GetFileURL returns the appropriate URL for a file. For S3 backends this is a
// presigned URL; for local storage it's the relative API path.
func (s *FileService) GetFileURL(ctx context.Context, file *models.File) (url string, expires *time.Time, err error) {
	if ps, ok := s.files.(storage.PresignCapable); ok {
		presigned, err := ps.PresignGet(ctx, file.StorageKey, presignTTL)
		if err != nil {
			return "", nil, fmt.Errorf("failed to generate presigned URL: %w", err)
		}
		exp := time.Now().Add(presignTTL)
		return presigned, &exp, nil
	}

	return fmt.Sprintf("/api/files/%s", file.ID), nil, nil
}

// GetFilesByNodeIDs batch-fetches files for the given node IDs and returns a
// map of nodeID → File. Used for enriching graph responses.
func (s *FileService) GetFilesByNodeIDs(ctx context.Context, nodeIDs []string) (map[string]models.File, error) {
	return s.repo.GetFilesByNodeIDs(ctx, nodeIDs)
}

// BuildFileInfo converts a models.File into a client-safe FileInfo with the
// correct URL for the current storage backend.
func (s *FileService) BuildFileInfo(ctx context.Context, file *models.File) (*models.FileInfo, error) {
	url, expires, err := s.GetFileURL(ctx, file)
	if err != nil {
		return nil, err
	}

	return &models.FileInfo{
		ID:         file.ID,
		Filename:   file.Filename,
		MimeType:   file.MimeType,
		SizeBytes:  file.SizeBytes,
		URL:        url,
		URLExpires: expires,
	}, nil
}

// IsPresignCapable reports whether the underlying storage supports presigned URLs.
func (s *FileService) IsPresignCapable() bool {
	_, ok := s.files.(storage.PresignCapable)
	return ok
}

// AttachFileMetadataToNodes batch-fetches files for reward nodes and attaches FileInfo.
func (s *FileService) AttachFileMetadataToNodes(ctx context.Context, nodes []models.Node) error {
	// Collect reward node IDs.
	var rewardNodeIDs []string
	for _, n := range nodes {
		if n.Type == models.RewardNode {
			rewardNodeIDs = append(rewardNodeIDs, n.ID)
		}
	}
	if len(rewardNodeIDs) == 0 {
		return nil
	}

	fileMap, err := s.GetFilesByNodeIDs(ctx, rewardNodeIDs)
	if err != nil {
		return err
	}

	for i := range nodes {
		if nodes[i].Type != models.RewardNode {
			continue
		}
		file, ok := fileMap[nodes[i].ID]
		if !ok {
			continue
		}
		rd, ok := nodes[i].Data.(*models.RewardData)
		if !ok {
			continue
		}
		info, err := s.BuildFileInfo(ctx, &file)
		if err != nil {
			return fmt.Errorf("failed to build file info for node %s: %w", nodes[i].ID, err)
		}
		rd.File = info
	}

	return nil
}

// PresignURL returns a presigned URL for the given file. Only valid if IsPresignCapable is true.
func (s *FileService) PresignURL(ctx context.Context, storageKey string) (string, error) {
	ps, ok := s.files.(storage.PresignCapable)
	if !ok {
		return "", fmt.Errorf("storage backend does not support presigned URLs")
	}
	return ps.PresignGet(ctx, storageKey, presignTTL)
}
