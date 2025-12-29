package sqlite

import (
	"backend/internal/models"
	"context"
	"strconv"
)

func (s *sqliteDB) CreateFavourRequest(ctx context.Context, request *models.FavourRequest) error {
	result, err := s.executor().ExecContext(ctx, `
		INSERT INTO favour_requests (favour_choice_id, message)
		VALUES (?, ?)
	`, request.ChoiceID, request.Message)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	request.ID = strconv.FormatInt(id, 10)
	return nil
}

func (s *sqliteDB) ListFavourRequests(ctx context.Context) ([]models.FavourRequest, error) {
	rows, err := s.executor().QueryContext(ctx, `
		SELECT id, favour_choice_id, message, requested_at, fulfilled_at
		FROM favour_requests
		ORDER BY requested_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]models.FavourRequest, 0)
	for rows.Next() {
		var request models.FavourRequest
		if err := rows.Scan(&request.ID, &request.ChoiceID, &request.Message, &request.RequestedAt, &request.FulfilledAt); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (s *sqliteDB) DeleteFavourRequest(ctx context.Context, requestID string) error {
	_, err := s.executor().ExecContext(ctx, `
		DELETE FROM favour_requests
		WHERE id = ?
	`, requestID)
	return err
}

func (s *sqliteDB) FulfilFavourRequest(ctx context.Context, requestID string) error {
	_, err := s.executor().ExecContext(ctx, `
		UPDATE favour_requests
		SET fulfilled_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, requestID)
	return err
}

func (s *sqliteDB) UnfulfilFavourRequest(ctx context.Context, requestID string) error {
	_, err := s.executor().ExecContext(ctx, `
		UPDATE favour_requests
		SET fulfilled_at = NULL
		WHERE id = ?
	`, requestID)
	return err
}
