package user

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maxpn01/x-twitter-clone/models"
	"github.com/maxpn01/x-twitter-clone/user/auth"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	usersByID          map[string]models.User
	usersByEmail       map[string]models.User
	usersByUsername    map[string]models.User
	refreshTokenUserID map[string]string
}

func (r *fakeUserRepository) CreateUser(ctx context.Context, input CreateUserInput) (models.User, error) {
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
func (r *fakeUserRepository) GetUserByID(ctx context.Context, id string) (models.User, error) {
	user, ok := r.usersByID[id]

	if !ok {
		return models.User{}, ErrUserNotFound
	}

	return user, nil
}
func (r *fakeUserRepository) GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (models.User, error) {
	if user, ok := r.usersByEmail[emailOrUsername]; ok {
		return user, nil
	}

	if user, ok := r.usersByUsername[emailOrUsername]; ok {
		return user, nil
	}

	return models.User{}, ErrUserNotFound
}
func (r *fakeUserRepository) StoreRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	r.refreshTokenUserID[refreshToken] = userID
	return nil
}
func (r *fakeUserRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	if _, ok := r.refreshTokenUserID[refreshToken]; !ok {
		return ErrRefreshTokenNotFound
	}

	delete(r.refreshTokenUserID, refreshToken)
	return nil
}

func TestUserService(t *testing.T) {
	t.Run("signup user successfully", func(t *testing.T) {
		service := newTestUserService(t)

		tokens, err := service.Signup(context.Background(), SignupInput{"test2@example.com", "testuser", "Test User", "Password123!"})
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

		createdUser, err := service.userRepo.GetUserByEmailOrUsername(context.Background(), "test2@example.com")
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

		_, err := service.Signup(context.Background(), SignupInput{"test@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.Signup(context.Background(), SignupInput{"test@example.com", "testuser2", "Test User", "Password123!"})
		if !errors.Is(err, ErrEmailAlreadyExists) {
			t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
		}
	})

	t.Run("signup rejects duplicate username", func(t *testing.T) {
		service := newTestUserService(t)

		_, err := service.Signup(context.Background(), SignupInput{"test@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.Signup(context.Background(), SignupInput{"test2@example.com", "testuser", "Test User", "Password123!"})
		if !errors.Is(err, ErrUsernameAlreadyExists) {
			t.Fatalf("expected ErrUsernameAlreadyExists, got %v", err)
		}
	})

	t.Run("signin user successfully", func(t *testing.T) {
		service := newTestUserService(t)

		_, err := service.Signup(context.Background(), SignupInput{"test@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		tokens, err := service.Signin(context.Background(), SigninInput{"test@example.com", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		if tokens.AccessToken == "" {
			t.Fatal("expected access token")
		}

		if tokens.RefreshToken == "" {
			t.Fatal("expected refresh token")
		}
	})

	t.Run("signin user successfully with username", func(t *testing.T) {
		service := newTestUserService(t)

		_, err := service.Signup(context.Background(), SignupInput{"test@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		tokens, err := service.Signin(context.Background(), SigninInput{"testuser", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		if tokens.AccessToken == "" {
			t.Fatal("expected access token")
		}

		if tokens.RefreshToken == "" {
			t.Fatal("expected refresh token")
		}
	})

	t.Run("signin requires email or username", func(t *testing.T) {
		service := newTestUserService(t)

		_, err := service.Signin(context.Background(), SigninInput{"", "Password123!"})
		if err == nil || err.Error() != "email_or_username is required" {
			t.Fatalf("expected email_or_username required error, got %v", err)
		}
	})

	t.Run("signout user successfully", func(t *testing.T) {
		service := newTestUserService(t)

		tokens, err := service.Signup(context.Background(), SignupInput{"test@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		err = service.Signout(context.Background(), SignoutInput{RefreshToken: tokens.RefreshToken})
		if err != nil {
			t.Fatal(err)
		}

		repo := service.userRepo.(*fakeUserRepository)
		if _, ok := repo.refreshTokenUserID[tokens.RefreshToken]; ok {
			t.Fatal("expected refresh token to be deleted")
		}
	})

	t.Run("signout rejects already revoked refresh token", func(t *testing.T) {
		service := newTestUserService(t)

		tokens, err := service.Signup(context.Background(), SignupInput{"test@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		err = service.Signout(context.Background(), SignoutInput{RefreshToken: tokens.RefreshToken})
		if err != nil {
			t.Fatal(err)
		}

		err = service.Signout(context.Background(), SignoutInput{RefreshToken: tokens.RefreshToken})
		if !errors.Is(err, ErrRefreshTokenNotFound) {
			t.Fatalf("expected ErrRefreshTokenNotFound, got %v", err)
		}
	})

	t.Run("signout rejects access token", func(t *testing.T) {
		service := newTestUserService(t)

		tokens, err := service.Signup(context.Background(), SignupInput{"test@example.com", "testuser", "Test User", "Password123!"})
		if err != nil {
			t.Fatal(err)
		}

		err = service.Signout(context.Background(), SignoutInput{RefreshToken: tokens.AccessToken})
		if err == nil {
			t.Fatal("expected access token to fail signout")
		}
	})
}

func TestUserHandler(t *testing.T) {
	t.Run("signup rejects incorrect json shape", func(t *testing.T) {
		service := newTestUserService(t)
		handler := NewUserHandler(service)

		body := `{
			"email": "test@example.com",
			"usernamea": "testuser",
			"fullname": "Test User",
			"password": "Password123!"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Signup(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		if strings.TrimSpace(rec.Body.String()) != "incorrect json shape for the request" {
			t.Fatalf("expected incorrect json shape error, got %q", rec.Body.String())
		}
	})
	t.Run("signin rejects incorrect json shape", func(t *testing.T) {
		service := newTestUserService(t)
		handler := NewUserHandler(service)

		body := `{
			"email": "test@example.com",
			"username": "testuser",
			"fullname": "Test User",
			"password": "Password123!"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Signin(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		if strings.TrimSpace(rec.Body.String()) != "incorrect json shape for the request" {
			t.Fatalf("expected incorrect json shape error, got %q", rec.Body.String())
		}
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
		usersByID:          map[string]models.User{},
		usersByEmail:       map[string]models.User{},
		usersByUsername:    map[string]models.User{},
		refreshTokenUserID: map[string]string{},
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
