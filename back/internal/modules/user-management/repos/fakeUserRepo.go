package user_repos

import (
	user_management_domain "chat/internal/modules/user-management/domain"
	"context"
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
		if user.GetEmail() == email {
			return user, nil
		}
	}
	return nil, nil
}
