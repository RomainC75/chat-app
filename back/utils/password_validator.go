package utils

import (
	"errors"
	"regexp"
)

func PasswordConstrainsValidator(password string) error{
	r, _ := regexp.Compile(`[\[\]*+,-.\/:;()<=>?@]`)
	specialCharacterResult := r.FindStringIndex(password)
	if len(specialCharacterResult)== 0 {
		return errors.New("need at least a special character in the password")
	}

	r, _ = regexp.Compile(`[A-Z]`)
	upperCharacterResult := r.FindStringIndex(password)
	if len(upperCharacterResult)== 0 {
		return errors.New("need at least an upper case in the password")
	}

	r, _ = regexp.Compile(`[0-9]`)
	digitCharacterResult := r.FindStringIndex(password)
	if len(digitCharacterResult)== 0 {
		return errors.New("need at least a digit in the password")
	}
	return nil
}