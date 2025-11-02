package user_management_domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	iD        uuid.UUID
	email     string
	password  string
	createdAt time.Time
	updatedAt time.Time
}

type UserSnapshot struct {
	ID        uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id uuid.UUID, email string, password string, now time.Time) *User {
	return &User{
		iD:        id,
		email:     email,
		password:  password,
		createdAt: now,
	}
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetID() uuid.UUID {
	return u.iD
}

func (u *User) GetToken(generatorFn func(user *User) (string, error)) (string, error) {
	return generatorFn(u)
}

func (u *User) IsAuthenticated(fn func(hash []byte, password []byte) error, hashPass string) error {
	return fn([]byte(u.password), []byte(hashPass))
}

func (u *User) ToSnapshot() UserSnapshot {
	return UserSnapshot{
		ID:        u.iD,
		Email:     u.email,
		CreatedAt: u.createdAt,
		Password:  u.password,
		UpdatedAt: u.updatedAt,
	}
}

func FromSnapshot(snapshot UserSnapshot) *User {
	return &User{
		iD:        snapshot.ID,
		email:     snapshot.Email,
		password:  snapshot.Password,
		createdAt: snapshot.CreatedAt,
		updatedAt: snapshot.UpdatedAt,
	}
}
