package user_management_app

import (
	"chat/internal/api/dto/requests"
	custom_errors "chat/internal/api/errors"
	"chat/internal/api/shared"
	user_management_domain "chat/internal/modules/user-management/domain"
	"chat/utils/encrypt"

	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type LogResponse struct {
	Id    int32  `json:"id"`
	Token string `json:"token"`
}

type BasicUserResponse struct {
	Id    int32  `json:"id"`
	Email string `json:"email"`
}

type VerifyResponse struct {
	BasicUserResponse
}

type CreateUserResponse struct {
	BasicUserResponse
}

type UserSrv struct {
	userRepo      user_management_domain.IUsers
	uuidGenerator shared.UUIDGenerator
	clock         shared.Clock
}

func NewUserSrv(userRepo user_management_domain.IUsers, uuidGenerator shared.UUIDGenerator, clock shared.Clock) *UserSrv {
	return &UserSrv{
		userRepo:      userRepo,
		clock:         clock,
		uuidGenerator: uuidGenerator,
	}
}

func (userSrv *UserSrv) CreateUserSrv(ctx context.Context, user requests.SignupRequest) (CreateUserResponse, error) {
	_, err := userSrv.userRepo.FindUserByEmail(ctx, user.Email)
	if err != nil && errors.As(err, &custom_errors.ErrNotFound{}) {
		return CreateUserResponse{}, errors.New("email already used")
	}

	b, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return CreateUserResponse{}, errors.New("error trying to encrypt the password")
	}

	newUser := user_management_domain.NewUser(
		userSrv.uuidGenerator.Generate(),
		user.Email,
		string(b),
		userSrv.clock.Now(),
	)

	err = userSrv.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return CreateUserResponse{}, err
	}
	return CreateUserResponse{
		BasicUserResponse: BasicUserResponse{
			Id:    int32(newUser.GetID().ID()),
			Email: newUser.GetEmail(),
		},
	}, nil
}

func (userSrv *UserSrv) LogUserSrv(ctx context.Context, user requests.LoginRequest) (LogResponse, error) {
	foundUser, err := userSrv.userRepo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return LogResponse{}, errors.New("user not found")
	}

	err = foundUser.IsAuthenticated(bcrypt.CompareHashAndPassword, user.Password)
	if err != nil {
		return LogResponse{}, errors.New("wrong password")
	}

	token, err := foundUser.GetToken(encrypt.Generate)
	if err != nil {
		return LogResponse{}, err
	}

	return LogResponse{
		Id:    int32(foundUser.GetID().ID()),
		Token: token,
	}, nil
}
