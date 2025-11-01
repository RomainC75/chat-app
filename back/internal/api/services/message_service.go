package services

import (
	db "chat/db/sqlc"
	"chat/internal/api/dto/requests"
	custom_errors "chat/internal/api/errors"
	repositories "chat/internal/api/repos"
	"chat/utils/encrypt"

	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type MessageSrv struct {
	userRepo repositories.UserRepositoryInterface
}

func NewMessageSrv() *MessageSrv {
	return &MessageSrv{
		userRepo: repositories.NewUserRepo(),
	}
}

func (messageSrv *MessageSrv) CreatemessageSrv(ctx context.Context, user requests.SignupRequest) (db.User, error) {
	_, err := messageSrv.userRepo.FindUserByEmail(ctx, user.Email)
	if err != nil && errors.As(err, &custom_errors.ErrNotFound{}) {
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

	createdUser, err := messageSrv.userRepo.CreateUser(ctx, userParams)
	if err != nil {
		return db.User{}, err
	}
	return createdUser, nil
}

func (messageSrv *MessageSrv) LogUserSrv(ctx context.Context, user requests.LoginRequest) (LogResponse, error) {
	foundUser, err := messageSrv.userRepo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return LogResponse{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		return LogResponse{}, errors.New("wrong password")
	}

	token, err := encrypt.Generate(foundUser)
	if err != nil {
		return LogResponse{}, err
	}

	return LogResponse{
		Id:    foundUser.ID,
		Token: token,
	}, nil
}
