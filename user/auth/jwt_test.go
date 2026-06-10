package auth

import (
	"testing"
	"time"
)

func TestJWTService(t *testing.T) {
	t.Run("generate and verify access token", func(t *testing.T) {
		service := newTestJWTService(t)

		pair, err := service.GenerateTokenPair("user-123")
		if err != nil {
			t.Fatal(err)
		}

		userID, err := service.VerifyToken(pair.AccessToken, AccessToken)
		if err != nil {
			t.Fatal(err)
		}

		if userID != "user-123" {
			t.Fatalf("got %q, want %q", userID, "user-123")
		}
	})

	t.Run("rejects wrong token type", func(t *testing.T) {
		service := newTestJWTService(t)

		pair, err := service.GenerateTokenPair("user-123")
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.VerifyToken(pair.RefreshToken, AccessToken)
		if err == nil {
			t.Fatal("expected refresh >token to fail as access token")
		}
	})

	t.Run("refresh access token", func(t *testing.T) {
		service := newTestJWTService(t)

		pair, err := service.GenerateTokenPair("user-123")
		if err != nil {
			t.Fatal(err)
		}

		accessToken, err := service.RefreshAccessToken(pair.RefreshToken)
		if err != nil {
			t.Fatal(err)
		}

		userID, err := service.VerifyToken(accessToken, AccessToken)
		if err != nil {
			t.Fatal(err)
		}

		if userID != "user-123" {
			t.Fatalf("got %q, want %q", userID, "user-123")
		}
	})
}

func newTestJWTService(t *testing.T) *JWTService {
	t.Helper()

	service, err := NewJWTService("secret", time.Minute, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	return service
}
