package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

func (s *sqliteDB) scanFile(scanner interface{ Scan(dest ...any) error }) (models.File, error) {
	var file models.File
	var fileID int64

	err := scanner.Scan(&fileID, &file.NodeID, &file.Filename, &file.MimeType, &file.StorageKey, &file.SizeBytes)
	if err != nil {
		return file, err
	}

	file.ID = fmt.Sprintf("%d", fileID)
	return file, nil
}

func (s *sqliteDB) GetFile(ctx context.Context, fileID string) (*models.File, error) {
	row := s.executor().QueryRowContext(ctx, `
		SELECT id, node_id, filename, mime_type, storage_key, size_bytes
		FROM reward_file
		WHERE id = ?
	`, fileID)

	file, err := s.scanFile(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get file %s: %w", fileID, err)
	}
	return &file, nil
}

func (s *sqliteDB) GetFileByNodeID(ctx context.Context, nodeID string) (*models.File, error) {
	row := s.executor().QueryRowContext(ctx, `
		SELECT id, node_id, filename, mime_type, storage_key, size_bytes
		FROM reward_file
		WHERE node_id = ?
	`, nodeID)

	file, err := s.scanFile(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get file for node %s: %w", nodeID, err)
	}
	return &file, nil
}

func (s *sqliteDB) GetFilesByNodeIDs(ctx context.Context, nodeIDs []string) (map[string]models.File, error) {
	if len(nodeIDs) == 0 {
		return nil, nil
	}

	jsonIDs, err := json.Marshal(nodeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal node IDs: %w", err)
	}

	rows, err := s.executor().QueryContext(ctx, `
		SELECT id, node_id, filename, mime_type, storage_key, size_bytes
		FROM reward_file
		WHERE node_id IN (SELECT value FROM json_each(?))
	`, string(jsonIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to query files by node IDs: %w", err)
	}
	defer rows.Close()

	result := make(map[string]models.File)
	for rows.Next() {
		file, err := s.scanFile(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file row: %w", err)
		}
		result[file.NodeID] = file
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating file rows: %w", err)
	}

	return result, nil
}

func (s *sqliteDB) DeleteFilesByNodeID(ctx context.Context, nodeID string) ([]string, error) {
	// First get the storage keys
	rows, err := s.executor().QueryContext(ctx, `
		SELECT storage_key FROM reward_file WHERE node_id = ?
	`, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query files for deletion: %w", err)
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("failed to scan storage key: %w", err)
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating storage keys: %w", err)
	}

	// Then delete the rows
	_, err = s.executor().ExecContext(ctx, `
		DELETE FROM reward_file WHERE node_id = ?
	`, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete files for node %s: %w", nodeID, err)
	}

	return keys, nil
}
