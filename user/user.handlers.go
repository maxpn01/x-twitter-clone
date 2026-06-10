package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	req := SignupInput{}

	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, "incorrect json shape for the request", http.StatusBadRequest)
		return
	}

	tokens, err := h.userService.Signup(r.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrEmailAlreadyExists) || errors.Is(err, ErrUsernameAlreadyExists) {
			status = http.StatusConflict
		}

		http.Error(w, err.Error(), status)
		return
	}

	if err := writeJSON(w, http.StatusOK, tokens); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *UserHandler) Signin(w http.ResponseWriter, r *http.Request) {
	req := SigninInput{}

	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, "incorrect json shape for the request", http.StatusBadRequest)
		return
	}

	tokens, err := h.userService.Signin(r.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrInvalidCredentials) {
			status = http.StatusUnauthorized
		}

		http.Error(w, err.Error(), status)
		return
	}

	if err := writeJSON(w, http.StatusOK, tokens); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *UserHandler) Signout(w http.ResponseWriter, r *http.Request) {
	refreshToken, ok := bearerToken(r)
	if !ok {
		http.Error(w, "the user is not authorized", http.StatusUnauthorized)
		return
	}

	if err := h.userService.Signout(r.Context(), SignoutInput{RefreshToken: refreshToken}); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func bearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}

	token, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok || token == "" {
		return "", false
	}

	return token, true
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, value any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes())
	return err
}
