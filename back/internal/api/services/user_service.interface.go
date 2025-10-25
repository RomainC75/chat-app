package services

import (
	db "chat/db/sqlc"
	"chat/internal/api/dto/requests"
	"context"
)

type UserServiceInterface interface {
	CreateUserSrv(ctx context.Context, user requests.SignupRequest) (db.User, error)
}
