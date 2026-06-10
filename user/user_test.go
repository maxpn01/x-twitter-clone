package user

import (
	"errors"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/maxpn01/x-twitter-clone/models"
	"github.com/maxpn01/x-twitter-clone/user/auth"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	usersByID       map[string]models.User
	usersByEmail    map[string]models.User
	usersByUsername map[string]models.User
}

func (r *fakeUserRepository) CreateUser(input CreateUserInput) (models.User, error) {
	if _, ok := r.usersByEmail[input.Email]; ok {
		return models.User{}, ErrEmailAlreadyExists
	}
	if _, ok := r.usersByUsername[input.Username]; ok {
		return models.User{}, ErrUsernameAlreadyExists
	}

	newUser := createNewTestUser(CreateUserInput{input.Email, input.Username, input.Fullname, input.PasswordHash})

	r.usersByID[newUser.ID] = newUser
	r.usersByEmail[newUser.Email] = newUser
	r.usersByUsername[newUser.Username] = newUser

	return newUser, nil
}
func (r *fakeUserRepository) GetUserByID(id string) (models.User, error) {
	user, ok := r.usersByID[id]

	if !ok {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}
func (r *fakeUserRepository) GetUserByEmail(email string) (models.User, error) {
	user, ok := r.usersByEmail[email]

	if !ok {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func TestUserService(t *testing.T) {
	t.Run("signup user successfully", func(t *testing.T) {
		service := newTestUserService(t)

		tokens, err := service.Signup(SignupInput{"test2@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		if tokens.AccessToken == "" {
			t.Fatal("expected access token")
		}

		if tokens.RefreshToken == "" {
			t.Fatal("expected refresh token")
		}

		userID, err := service.jwtService.VerifyToken(tokens.AccessToken, auth.AccessToken)
		if err != nil {
			t.Fatal(err)
		}

		if userID == "" {
			t.Fatal("expected user id in token subject")
		}

		createdUser, err := service.userRepo.GetUserByEmail("test2@example.com")
		if err != nil {
			t.Fatal(err)
		}

		if createdUser.PasswordHash == "Password123!" {
			t.Fatal("expected password to be hashed")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(createdUser.PasswordHash), []byte("Password123!")); err != nil {
			t.Fatal("expected password hash to match password")
		}
	})

	t.Run("signup rejects duplicate email", func(t *testing.T) {
		service := newTestUserService(t)

		_, err := service.Signup(SignupInput{"test@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.Signup(SignupInput{"test@example.com", "testuser2", "Test User", "Password123!"})
		if !errors.Is(err, ErrEmailAlreadyExists) {
			t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
		}
	})

	t.Run("signup rejects duplicate username", func(t *testing.T) {
		service := newTestUserService(t)

		_, err := service.Signup(SignupInput{"test@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.Signup(SignupInput{"test2@example.com", "testuser", "Test User", "Password123!"})
		if !errors.Is(err, ErrUsernameAlreadyExists) {
			t.Fatalf("expected ErrUsernameAlreadyExists, got %v", err)
		}
	})

	t.Run("signin user successfully", func(t *testing.T) {
		t.Fatal("not implemented")
	})

	t.Run("signout user successfully", func(t *testing.T) {
		t.Fatal("not implemented")
	})
}

func newTestUserService(t *testing.T) *UserService {
	t.Helper()

	userRepo := newTestFakeUserRepo(t)
	jwtService := newTestJWTService(t)

	service, err := NewUserService(userRepo, jwtService)
	if err != nil {
		t.Fatal(err)
	}

	return service
}

func newTestFakeUserRepo(t *testing.T) *fakeUserRepository {
	t.Helper()

	return &fakeUserRepository{
		usersByID:       map[string]models.User{},
		usersByEmail:    map[string]models.User{},
		usersByUsername: map[string]models.User{},
	}
}

func newTestJWTService(t *testing.T) *auth.JWTService {
	t.Helper()

	service, err := auth.NewJWTService("secret", time.Minute, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	return service
}

func createNewTestUser(input CreateUserInput) models.User {
	return models.User{
		ID:           strconv.Itoa(rand.Intn(10)),
		Email:        input.Email,
		Username:     input.Username,
		Fullname:     input.Fullname,
		PasswordHash: input.PasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
