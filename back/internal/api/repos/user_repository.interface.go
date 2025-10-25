package repositories

import (
	db "chat/db/sqlc"
	"context"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	FindUserByEmail(ctx context.Context, email string) (db.User, error)
}
