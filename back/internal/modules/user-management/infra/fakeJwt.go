package user_management_infra

import (
	user_management_domain "chat/internal/modules/user-management/domain"
	user_management_encrypt "chat/internal/modules/user-management/domain/encrypt"
)

type FakeJWT struct {
	ExpectedJwt string
}

func NewFakeJWT() *FakeJWT {
	return &FakeJWT{}
}

func (f *FakeJWT) Generate(user *user_management_domain.User) (string, error) {
	return f.ExpectedJwt, nil
}

func (f *FakeJWT) GetClaimsFromToken(tokenString string) (user_management_encrypt.JwtClaim, error) {
	return user_management_encrypt.JwtClaim{}, nil
}
