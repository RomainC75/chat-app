package encrypt

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, receivedPwd string) error{
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(receivedPwd))
	if err != nil {
		return err
	}
	return nil
}
