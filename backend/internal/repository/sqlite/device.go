package sqlite

import (
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
