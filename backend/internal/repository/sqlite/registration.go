package sqlite

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
)

// RegistrationId is the fixed ID used to store the registration code.
// There is only one registration code at any time.
const RegistrationId = 1

func (s *sqliteDB) RefreshRegistrationCode(ctx context.Context, model models.RegistrationCode) error {
	_, err := s.db.ExecContext(ctx, `
			INSERT OR REPLACE INTO registration_codes (id, code, expires_at)
			VALUES (?, ?, ?);

	`, RegistrationId, model.Code, model.ExpiresAt)
	return err
}

func (s *sqliteDB) DeleteRegistrationCode(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
			DELETE FROM registration_codes WHERE id = ?;
	`, RegistrationId)
	return err
}

func (s *sqliteDB) GetRegistrationCode(ctx context.Context) (model *models.RegistrationCode, err error) {
	row := s.db.QueryRowContext(ctx, `
			SELECT code, expires_at FROM registration_codes WHERE id = ?;
	`, RegistrationId)

	model = &models.RegistrationCode{}
	err = row.Scan(&model.Code, &model.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return model, nil
}
