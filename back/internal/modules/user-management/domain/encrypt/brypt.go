package user_management_encrypt

type Bcrypt interface {
	HashAndSalt(pwd string) ([]byte, error)
	ComparePasswords(hashedPwd string, receivedPwd string) error
}
