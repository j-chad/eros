package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

func (s *sqliteDB) RegisterDevice(ctx context.Context, token string, deviceInfo string, expiry time.Time) error {
	_, err := s.executor().ExecContext(ctx, `
		INSERT INTO device (token, device_info, expires_at)
		VALUES (?, ?, ?)
	`, token, deviceInfo, expiry)
	return err
}

func (s *sqliteDB) GetDeviceExpiryByToken(ctx context.Context, token string) (expiresAt *time.Time, err error) {
	row := s.executor().QueryRowContext(ctx, `
		SELECT expires_at
		FROM device
		WHERE token = ? AND expires_at > CURRENT_TIMESTAMP
	`, token)
	err = row.Scan(&expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return expiresAt, err
}

func (s *sqliteDB) ListDevices(ctx context.Context) ([]models.Device, error) {
	rows, err := s.executor().QueryContext(ctx, `
		SELECT id, device_info, expires_at, registered_at, last_seen_at
		FROM device
		ORDER BY registered_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		if err := rows.Scan(&device.ID, &device.DeviceInfo, &device.ExpiresAt, &device.RegisteredAt, &device.LastSeenAt); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, rows.Err()
}

func (s *sqliteDB) DeleteDevice(ctx context.Context, deviceID string) error {
	_, err := s.executor().ExecContext(ctx, `
		DELETE FROM device
		WHERE id = ?
	`, deviceID)
	return err
}
