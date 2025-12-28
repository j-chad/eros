package sqlite

import (
	"context"
	"time"
)

func (s *sqliteDB) RegisterDevice(ctx context.Context, token string, deviceInfo string, expiry time.Time) error {
	_, err := s.executor().ExecContext(ctx, `
		INSERT INTO device (token, device_info, expires_at)
		VALUES (?, ?, ?)
	`, token, deviceInfo, expiry)
	return err
}
