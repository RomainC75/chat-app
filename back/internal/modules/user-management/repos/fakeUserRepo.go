package user_repos

import (
	custom_errors "chat/internal/api/errors"
	user_management_domain "chat/internal/modules/user-management/domain"
	"context"
	"errors"
	"fmt"
)

type FakeUserRepo struct {
	SavedUser *user_management_domain.User
	Users     []*user_management_domain.User
}

func NewFakeUserRepo() *FakeUserRepo {
	return &FakeUserRepo{}
}

func (fur *FakeUserRepo) CreateUser(ctx context.Context, newUser *user_management_domain.User) error {
	fur.SavedUser = newUser
	return nil
}

func (fur *FakeUserRepo) FindUserByEmail(ctx context.Context, email string) (*user_management_domain.User, error) {
	for _, user := range fur.Users {
		fmt.Println(email, user.ToSnapshot())
		if user.GetEmail() == email {
			return user, nil
		}
	}
	return nil, custom_errors.NewErrNotFound(errors.New("user not found"))
}
