package user_management_encrypt

type Bcrypt interface {
	HashAndSalt(pwd string) (string, error)
	ComparePasswords(hashedPwd string, receivedPwd string) error
}
