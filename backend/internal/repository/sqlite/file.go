package sqlite

import (
	"backend/internal/models"
	"context"
	"fmt"
)

func (s *sqliteDB) CreateFile(ctx context.Context, file *models.File) error {
	result, err := s.executor().ExecContext(ctx, `
		INSERT INTO reward_file (node_id, filename, mime_type, storage_key, size_bytes)
		VALUES (?, ?, ?, ?, ?)
	`, file.NodeID, file.Filename, file.MimeType, file.StorageKey, file.SizeBytes)
	if err != nil {
		return fmt.Errorf("failed to insert file row: %w", err)
	}

	fileID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get ID for file: %w", err)
	}

	file.ID = fmt.Sprintf("%d", fileID)
	return nil
}

func (s *sqliteDB) ListFiles(ctx context.Context, nodeID string) ([]models.File, error) {
	rows, err := s.executor().QueryContext(ctx, `
		SELECT id, node_id, filename, mime_type, storage_key, size_bytes
		FROM reward_file
		WHERE node_id = ?
	`, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query files for node %s: %w", nodeID, err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		var fileID int64

		if err := rows.Scan(&fileID, &file.NodeID, &file.Filename, &file.MimeType, &file.StorageKey, &file.SizeBytes); err != nil {
			return nil, fmt.Errorf("failed to scan file row: %w", err)
		}

		file.ID = fmt.Sprintf("%d", fileID)
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating file rows: %w", err)
	}

	return files, nil
}
