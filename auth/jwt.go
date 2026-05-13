package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken(sub, email, username string, expiry time.Duration) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"sub":      sub,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(secretKey))

	return signedString, err
}

func GenerateAccessToken(id, email, username string) (string, error) {
	return generateToken(id, email, username, 15*time.Minute)
}

func GenerateRefreshToken(id, email, username string) (string, error) {
	return generateToken(id, email, username, 7*24*time.Hour)
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
