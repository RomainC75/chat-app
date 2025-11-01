package repositories

import (
	db "chat/db/sqlc"
	custom_errors "chat/internal/api/errors"
	"context"
	"database/sql"
	"fmt"

	"time"
)

type RoomRepository struct {
	Store *db.Store
}

func NewRoomRepo() *RoomRepository {
	return &RoomRepository{
		Store: db.GetConnection(),
	}
}

func (msgRepo *RoomRepository) CreateRoom(ctx context.Context, arg db.CreateUserParams) (db.Message, error) {
	arg.CreatedAt = time.Now()
	arg.UpdatedAt = arg.CreatedAt
	message, err := (*msgRepo.Store).CreateUser(ctx, arg)
	if err == sql.ErrNoRows {
		return db.Message{}, custom_errors.NewErrNotFound(err)
	} else if err != nil {
		return db.Message{}, err
	}
	return message, nil
}

func (userRepo *RoomRepository) FindUserByEmail(ctx context.Context, email string) (db.User, error) {
	foundUser, err := (*userRepo.Store).GetUserByEmail(ctx, email)
	if err != nil {
		fmt.Println("find", err)
		return db.User{}, err
	}
	return foundUser, nil
}
