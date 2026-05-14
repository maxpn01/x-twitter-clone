package auth

import (
	"errors"

	"github.com/maxpn01/x-twitter-clone/models"
	"github.com/maxpn01/x-twitter-clone/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Signup(email, username, password string) (AuthTokens, error) {
	email = normalizeAuthInput(email)
	username = normalizeAuthInput(username)

	if err := validateEmail(email); err != nil {
		return AuthTokens{}, err
	}
	if err := validateUsername(username); err != nil {
		return AuthTokens{}, err
	}
	if err := validatePassword(password, email, username); err != nil {
		return AuthTokens{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthTokens{}, err
	}

	user := models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
	}

	user, err = s.userRepo.CreateUser(user)
	if err != nil {
		return AuthTokens{}, err
	}
	if user.ID == "" {
		return AuthTokens{}, errors.New("created user is missing id")
	}

	accessToken, err := GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		return AuthTokens{}, err
	}
	refreshToken, err := GenerateRefreshToken(user.ID, user.Email, user.Username)
	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
