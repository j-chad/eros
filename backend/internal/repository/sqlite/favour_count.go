package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
)

func (s *sqliteDB) GetFavourCount(ctx context.Context) (models.FavourCount, error) {
	var count models.FavourCount
	err := s.executor().QueryRowContext(ctx, `
		SELECT
			fc.total_favours AS total_favours,
			fc.total_favours - COUNT(fr.id) AS remaining_favours
		FROM favour_count fc
		LEFT JOIN favour_requests fr
			ON 1 = 1
		WHERE fc.id = 0;
	`).Scan(&count.Total, &count.Remaining)

	if errors.Is(err, sql.ErrNoRows) {
		return models.FavourCount{Total: 0, Remaining: 0}, nil
	}

	return count, err
}

func (s *sqliteDB) UpdateFavourCount(ctx context.Context, count int) error {
	_, err := s.executor().ExecContext(ctx, `INSERT INTO favour_count (id, total_favours) VALUES (0, ?) ON CONFLICT DO UPDATE SET total_favours=excluded.total_favours`, count)
	return err
}
