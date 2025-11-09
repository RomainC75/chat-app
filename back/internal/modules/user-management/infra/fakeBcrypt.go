package user_management_infra

import "errors"

type FakeBcrypt struct {
	ExpectedHash          string
	ExpectedCompareResult bool
}

func NewFakeBcrypt() *FakeBcrypt {
	return &FakeBcrypt{}
}

func (fb *FakeBcrypt) HashAndSalt(pwd string) ([]byte, error) {
	return []byte(pwd), nil
}

func (fb *FakeBcrypt) ComparePasswords(hashedPwd string, receivedPwd string) error {
	if fb.ExpectedCompareResult {
		return nil
	} else {
		return errors.New("passwords do not match")
	}
}
