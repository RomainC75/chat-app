package repositories

import (
	db "chat/db/sqlc"
	"context"

	"time"
)

type UserRepository struct {
	Store *db.Store
}

func NewUserRepo() *UserRepository {
	return &UserRepository{
		Store: db.GetConnection(),
	}
}

func (userRepo *UserRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {

	arg.CreatedAt = time.Now()
	arg.UpdatedAt = arg.CreatedAt
	user, err := (*userRepo.Store).CreateUser(ctx, arg)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (userRepo *UserRepository) FindUserByEmail(ctx context.Context, email string) (db.User, error) {
	foundUser, err := (*userRepo.Store).GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}
	return foundUser, nil
}
