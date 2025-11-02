package user_management_encrypt

import (
	user_management_domain "chat/internal/modules/user-management/domain"
)

type JwtClaim map[string]interface{}

type JWT interface {
	Generate(user *user_management_domain.User) (string, error)
	GetClaimsFromToken(tokenString string) (JwtClaim, error)
}
