package common

import "golang.org/x/crypto/bcrypt"

func HassPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

func CheckPassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword),[]byte(password))
}