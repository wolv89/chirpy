package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {

	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", err
	}

	return string(pass), nil

}

func CheckPasswordHash(password, hash string) error {

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

}
