package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(incomingPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(incomingPassword)); err != nil {
		return err
	}
	return nil
}
