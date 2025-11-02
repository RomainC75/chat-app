package user_management_domain

import (
	"context"
)

type IUsers interface {
	CreateUser(ctx context.Context, newUser *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}
