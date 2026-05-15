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

func (s *sqliteDB) GetDeviceByToken(ctx context.Context, token string) (*models.Device, error) {
	var device models.Device

	//language=sqlite
	row := s.executor().QueryRowContext(ctx, `
		SELECT id, device_info, registered_at, last_seen_at, expires_at
		FROM device
		WHERE token = ? AND expires_at > CURRENT_TIMESTAMP
	`, token)
	err := row.Scan(&device.ID, &device.DeviceInfo, &device.RegisteredAt, &device.LastSeenAt, &device.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &device, err
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

	devices := make([]models.Device, 0)
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

func (s *sqliteDB) UpdateDeviceInfo(ctx context.Context, deviceID string, deviceInfo string) error {
	_, err := s.executor().ExecContext(ctx, `
		UPDATE device
		SET device_info = ?
		WHERE id = ?
	`, deviceInfo, deviceID)
	return err
}

func (s *sqliteDB) UpdateDeviceLastSeenByToken(ctx context.Context, token string, lastSeen time.Time) error {
	_, err := s.executor().ExecContext(ctx, `
		UPDATE device
		SET last_seen_at = ?
		WHERE token = ?
	`, lastSeen, token)
	return err
}
