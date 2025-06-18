package auth

import "context"

type Service interface {
	Login(ctx context.Context, username, password string) (string, error) // возвращает JWT
	Register(ctx context.Context, username, password string) error
	ValidateToken(ctx context.Context, token string) (userID string, err error)
}
