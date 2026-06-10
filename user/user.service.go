package user

import (
	"errors"

	"github.com/maxpn01/x-twitter-clone/user/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   UserRepository
	jwtService *auth.JWTService
}

func NewUserService(userRepo UserRepository, jwtService *auth.JWTService) (*UserService, error) {
	if userRepo == nil {
		return nil, errors.New("user repository is required")
	}
	if jwtService == nil {
		return nil, errors.New("jwt service is required")
	}

	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}, nil
}

type SignupInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}

func (s *UserService) Signup(input SignupInput) (auth.TokenPair, error) {
	err := ValidateSignupInput(input)
	if err != nil {
		return auth.TokenPair{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return auth.TokenPair{}, err
	}

	user, err := s.userRepo.CreateUser(CreateUserInput{
		input.Email,
		input.Username,
		input.Fullname,
		string(passwordHash),
	})
	if err != nil {
		return auth.TokenPair{}, err
	}

	tokens, err := s.jwtService.GenerateTokenPair(user.ID)
	if err != nil {
		return auth.TokenPair{}, err
	}

	if err := s.userRepo.StoreRefreshToken(user.ID, tokens.RefreshToken); err != nil {
		return auth.TokenPair{}, err
	}

	return tokens, nil
}

type SigninInput struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

func (s *UserService) Signin(input SigninInput) (auth.TokenPair, error) {
	if err := ValidateSigninInput(input); err != nil {
		return auth.TokenPair{}, err
	}

	user, err := s.userRepo.GetUserByEmailOrUsername(input.EmailOrUsername)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return auth.TokenPair{}, ErrInvalidCredentials
		}
		return auth.TokenPair{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return auth.TokenPair{}, ErrInvalidCredentials
	}

	tokens, err := s.jwtService.GenerateTokenPair(user.ID)
	if err != nil {
		return auth.TokenPair{}, err
	}

	if err := s.userRepo.StoreRefreshToken(user.ID, tokens.RefreshToken); err != nil {
		return auth.TokenPair{}, err
	}

	return tokens, nil
}

type SignoutInput struct {
	RefreshToken string
}

func (s *UserService) Signout(input SignoutInput) error {
	_, err := s.jwtService.VerifyToken(input.RefreshToken, auth.RefreshToken)
	if err != nil {
		return err
	}

	return s.userRepo.DeleteRefreshToken(input.RefreshToken)
}
