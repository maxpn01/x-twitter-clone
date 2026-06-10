package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	Type TokenType `json:"type"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTService struct {
	secretKey  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTService(secret string, accessTTL, refreshTTL time.Duration) (*JWTService, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if accessTTL <= 0 {
		return nil, errors.New("access token ttl must be positive")
	}
	if refreshTTL <= 0 {
		return nil, errors.New("refresh token ttl must be positive")
	}

	return &JWTService{
		secretKey:  []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

func (s *JWTService) IssueToken(sub string, tokenType TokenType, exp time.Time) (string, error) {
	tokenID, err := randomTokenID()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID,
		},
		Type: tokenType,
	})
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func randomTokenID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (s *JWTService) GenerateTokenPair(userID string) (TokenPair, error) {
	accessToken, err := s.IssueToken(userID, AccessToken, time.Now().Add(s.accessTTL))
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := s.IssueToken(userID, RefreshToken, time.Now().Add(s.refreshTTL))
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *JWTService) VerifyToken(tokenString string, tokenType TokenType) (string, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return s.secretKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if claims.Type != tokenType {
		return "", fmt.Errorf("invalid token type")
	}

	if claims.Subject == "" {
		return "", fmt.Errorf("missing token subject")
	}

	return claims.Subject, nil
}

func (s *JWTService) RefreshAccessToken(refreshToken string) (string, error) {
	sub, err := s.VerifyToken(refreshToken, RefreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err := s.IssueToken(sub, AccessToken, time.Now().Add(s.accessTTL))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
