package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

func (s *sqliteDB) GetFavourCount(ctx context.Context) (int, error) {
	var count int
	err := s.executor().QueryRowContext(ctx, `SELECT total_favours FROM favour_count LIMIT 1`).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (s *sqliteDB) UpdateFavourCount(ctx context.Context, count int) error {
	_, err := s.executor().ExecContext(ctx, `INSERT INTO favour_count (id, total_favours) VALUES (0, ?) ON CONFLICT DO UPDATE SET total_favours=excluded.total_favours`, count)
	return err
}
