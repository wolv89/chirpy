package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	now := time.Now()
	expiry := now.Add(expiresIn)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now.UTC()),
		ExpiresAt: jwt.NewNumericDate(expiry.UTC()),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return ss, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(id)

}

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	if !strings.Contains(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header")
	}

	return strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1)), nil

}

func MakeRefreshToken() string {

	token := make([]byte, 32)
	rand.Read(token)

	return hex.EncodeToString(token)

}

func GetAPIKey(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	if !strings.Contains(authHeader, "ApiKey ") {
		return "", errors.New("invalid authorization header")
	}

	return strings.TrimSpace(strings.Replace(authHeader, "ApiKey", "", 1)), nil

}
