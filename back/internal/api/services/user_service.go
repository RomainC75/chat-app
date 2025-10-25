package services

import (
	db "chat/db/sqlc"
	"chat/internal/api/dto/requests"
	repositories "chat/internal/api/repos"

	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserSrv struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserSrv() *UserSrv {
	return &UserSrv{
		userRepo: repositories.NewUserRepo(),
	}
}

func (userSrv *UserSrv) CreateUserSrv(ctx context.Context, user requests.SignupRequest) (db.User, error) {
	_, err := userSrv.userRepo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return db.User{}, errors.New("email already used")
	}

	b, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, errors.New("error trying to encrypt the password")
	}

	userParams := db.CreateUserParams{
		Email:    user.Email,
		Password: string(b),
	}

	createdUser, err := userSrv.userRepo.CreateUser(ctx, userParams)
	if err != nil {
		return db.User{}, err
	}
	return createdUser, nil
}
