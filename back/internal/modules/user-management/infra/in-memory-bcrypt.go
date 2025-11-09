package user_management_infra

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
}

func NewInMemoryBcrypt() *Bcrypt {
	return &Bcrypt{}
}

func (b *Bcrypt) HashAndSalt(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return hash, nil
}

func (b *Bcrypt) ComparePasswords(hashedPwd string, receivedPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(receivedPwd))
	if err != nil {
		return err
	}
	return nil
}
