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

type LogResponse struct {
	Id    int32  `json:"id"`
	Token string `json:"token"`
}

type VerifyResponse struct {
	Id    int32  `json:"id"`
	Email string `json:"email"`
}

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

	createdUser, err := userSrv.userRepo.CreateUser(ctx, userParams)
	if err != nil {
		return db.User{}, err
	}
	return createdUser, nil
}

func (userSrv *UserSrv) LogUserSrv(ctx context.Context, user requests.LoginRequest) (LogResponse, error) {
	foundUser, err := userSrv.userRepo.FindUserByEmail(ctx, user.Email)
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
