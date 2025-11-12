package user_management_infra

import (
	"chat/config"
	user_management_domain "chat/internal/modules/user-management/domain"
	user_management_encrypt "chat/internal/modules/user-management/domain/encrypt"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type InMemoryJWT struct {
}

func NewInMemoryJWT() *InMemoryJWT {
	return &InMemoryJWT{}
}

type Claims struct {
	*jwt.RegisteredClaims
	ID    uuid.UUID
	Email string
}

func (JWT *InMemoryJWT) Generate(user *user_management_domain.User) (string, error) {
	secret := config.Get().Jwt.Secret

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	exp := time.Now().Add(time.Hour * 24)

	token.Claims = &Claims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   user.GetID().String(),
		},
		user.GetID(),
		user.GetEmail(),
	}
	val, err := token.SignedString([]byte(secret))

	if err != nil {
		return "error trying to set the token", err
	}
	return val, nil
}

func (JWT *InMemoryJWT) GetClaimsFromToken(tokenString string) (user_management_encrypt.JwtClaim, error) {
	secret := config.Get().Jwt.Secret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return user_management_encrypt.JwtClaim(claims), nil
	}
	return nil, err
}
