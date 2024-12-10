package password

import (
	"golang.org/x/crypto/bcrypt"
)

func NewBcryptPasswordAlgorithm() PasswordAlgorithm {
	return bcryptPasswordAlgorithm{}
}

type bcryptPasswordAlgorithm struct {
}

func (bpa bcryptPasswordAlgorithm) HidePassword(password, salt string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(combinePasswordAndSalt(password, salt)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (bpa bcryptPasswordAlgorithm) CheckPassword(password, salt, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(combinePasswordAndSalt(password, salt)))
	return err == nil
}
