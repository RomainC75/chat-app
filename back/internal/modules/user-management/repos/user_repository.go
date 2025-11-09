package user_repos

import (
	db "chat/db/sqlc"
	custom_errors "chat/internal/api/errors"
	user_management_domain "chat/internal/modules/user-management/domain"
	"context"
	"database/sql"
	"fmt"

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

func (userRepo *UserRepository) CreateUser(ctx context.Context, newUser *user_management_domain.User) error {
	userSnapshot := newUser.ToSnapshot()
	arg := db.CreateUserParams{
		ID:        userSnapshot.ID,
		Email:     userSnapshot.Email,
		Password:  userSnapshot.Password,
		CreatedAt: userSnapshot.CreatedAt,
		UpdatedAt: userSnapshot.UpdatedAt,
	}
	arg.CreatedAt = time.Now()
	arg.UpdatedAt = arg.CreatedAt
	_, err := (*userRepo.Store).CreateUser(ctx, arg)
	if err == sql.ErrNoRows {
		return custom_errors.NewErrNotFound(err)
	} else if err != nil {
		return err
	}
	return nil
}

func (userRepo *UserRepository) FindUserByEmail(ctx context.Context, email string) (*user_management_domain.User, error) {
	foundUser, err := (*userRepo.Store).GetUserByEmail(ctx, email)
	if err != nil {
		fmt.Println("find", err)
		return nil, err
	}
	userSnapshot := user_management_domain.UserSnapshot{
		ID:        foundUser.ID,
		Email:     foundUser.Email,
		Password:  foundUser.Password,
		CreatedAt: foundUser.CreatedAt,
		UpdatedAt: foundUser.UpdatedAt,
	}
	return user_management_domain.FromSnapshot(userSnapshot), nil
}
