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

	return s.jwtService.GenerateTokenPair(user.ID)
}

// func (s *UserService) Signin(email, username, password string) (auth.TokenPair, error)

// func (s *UserService) Signout(notsure string) error
