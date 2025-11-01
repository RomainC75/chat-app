package repos

import (
	db "chat/db/sqlc"
	"context"
)

type MessageRepositoryInterface interface {
	CreateMessage(ctx context.Context, arg db.CreateMessageParams) (db.Message, error)
	// FindUserByEmail(ctx context.Context, email string) (db.User, error)
}
