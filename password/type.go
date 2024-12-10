package password

type PasswordAlgorithm interface {
	HidePassword(password, salt string) (string, error)
	CheckPassword(password, salt, hashedPassword string) bool
}

func combinePasswordAndSalt(password, salt string) string {
	return password + salt
}
