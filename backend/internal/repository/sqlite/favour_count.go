package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
)

func (s *sqliteDB) GetFavourCount(ctx context.Context) (models.FavourCount, error) {
	var total, remaining sql.NullInt64
	err := s.executor().QueryRowContext(ctx, `
		SELECT
			count.total_favours AS total_favours,
			count.total_favours - COALESCE(
				(
					SELECT SUM(choice.cost)
					FROM favour_requests requests
					JOIN favour_choice choice
					  ON choice.id = requests.favour_choice_id
				),
				0
			) AS remaining_favours
		FROM favour_count count
		WHERE count.id = 0;
	`).Scan(&total, &remaining)

	if errors.Is(err, sql.ErrNoRows) {
		return models.FavourCount{Total: 0, Remaining: 0}, nil
	}

	return models.FavourCount{
		Total:     int(total.Int64),
		Remaining: int(remaining.Int64),
	}, err
}

func (s *sqliteDB) UpdateFavourCount(ctx context.Context, count int) error {
	_, err := s.executor().ExecContext(ctx, `INSERT INTO favour_count (id, total_favours) VALUES (0, ?) ON CONFLICT DO UPDATE SET total_favours=excluded.total_favours`, count)
	return err
}
