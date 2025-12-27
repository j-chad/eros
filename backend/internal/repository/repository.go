package repository

import "context"

type Repository interface {
	Close()
	WithTx(ctx context.Context, fn func(Repository) error) error

	//RegistrationRepository
}

type RegistrationRepository interface {
	CreateRegistrationCode(ctx context.Context, token string) error
}
