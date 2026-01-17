package sqlite

import (
	"backend/internal/models"
	"context"
	"strconv"
)

// CreateFavourChoice inserts a new favour choice into the database.
// The ID of the newly created choice is set in the provided choice object.
func (s *sqliteDB) CreateFavourChoice(ctx context.Context, choice *models.FavourChoice) error {
	result, err := s.executor().ExecContext(ctx, `
		INSERT INTO favour_choice (label, description, cost, can_message)
		VALUES (?, ?, ?, ?)
	`, choice.Label, choice.Description, choice.Cost, choice.CanMessage)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	choice.ID = strconv.FormatInt(id, 10)
	return nil
}

func (s *sqliteDB) UpdateFavourChoice(ctx context.Context, choice models.FavourChoice) error {
	_, err := s.executor().ExecContext(ctx, `
		UPDATE favour_choice
		SET label = ?, description = ?, can_message = ?, cost = ?
		WHERE id = ?
	`, choice.Label, choice.Description, choice.CanMessage, choice.Cost, choice.ID)
	return err
}

func (s *sqliteDB) DeleteFavourChoice(ctx context.Context, choiceID string) error {
	_, err := s.executor().ExecContext(ctx, `
		DELETE FROM favour_choice
		WHERE id = ?
	`, choiceID)
	return err
}

func (s *sqliteDB) ListFavourChoices(ctx context.Context) ([]models.FavourChoice, error) {
	rows, err := s.executor().QueryContext(ctx, `
		SELECT id, label, description, cost, can_message, created_at, updated_at
		FROM favour_choice
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	choices := make([]models.FavourChoice, 0)
	for rows.Next() {
		var choice models.FavourChoice
		if err := rows.Scan(&choice.ID, &choice.Label, &choice.Description, &choice.Cost, &choice.CanMessage, &choice.CreatedAt, &choice.UpdatedAt); err != nil {
			return nil, err
		}
		choices = append(choices, choice)
	}

	return choices, rows.Err()
}

func (s *sqliteDB) GetFavourCostByID(ctx context.Context, choiceID string) (int, error) {
	var cost int
	err := s.executor().QueryRowContext(ctx, `
		SELECT cost
		FROM favour_choice
		WHERE id = ?
	`, choiceID).Scan(&cost)
	if err != nil {
		return 0, err
	}
	return cost, nil
}
