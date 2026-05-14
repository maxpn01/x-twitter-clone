package auth

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/maxpn01/x-twitter-clone/apperror"
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

func (s *AuthService) Signin(email, username, password string) (AuthTokens, error) {
	email = normalizeAuthInput(email)
	username = normalizeAuthInput(username)

	if password == "" {
		return AuthTokens{}, apperror.New(apperror.CodeValidation, "password is required")
	}

	var user models.User
	var err error
	switch {
	case email != "" && username != "":
		return AuthTokens{}, apperror.New(apperror.CodeValidation, "provide either email or username, not both")
	case email != "":
		if err := validateEmail(email); err != nil {
			return AuthTokens{}, err
		}
		user, err = s.userRepo.GetUserByEmail(email)
	case username != "":
		if err := validateUsername(username); err != nil {
			return AuthTokens{}, err
		}
		user, err = s.userRepo.GetUserByUsername(username)
	default:
		return AuthTokens{}, apperror.New(apperror.CodeValidation, "email or username is required")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return AuthTokens{}, invalidCredentialsError()
	}
	if err != nil {
		return AuthTokens{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return AuthTokens{}, invalidCredentialsError()
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

func (s *AuthService) Signout(accessToken string) error {
	if strings.TrimSpace(accessToken) == "" {
		return apperror.New(apperror.CodeUnauthorized, "bearer token is required")
	}
	if _, err := VerifyToken(accessToken); err != nil {
		return apperror.New(apperror.CodeUnauthorized, "invalid or expired token")
	}

	return nil
}

func invalidCredentialsError() error {
	return apperror.New(apperror.CodeInvalidCredentials, "invalid email, username or password")
}
