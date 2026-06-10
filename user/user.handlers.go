package user

import (
	"encoding/json"
	"errors"
	"net/http"
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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := h.userService.Signup(req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrEmailAlreadyExists) || errors.Is(err, ErrUsernameAlreadyExists) {
			status = http.StatusConflict
		}

		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *UserHandler) Signin(w http.ResponseWriter, r *http.Request) {
	req := SigninInput{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "incorrect json shape for the request", http.StatusBadRequest)
		return
	}

	tokens, err := h.userService.Signin(req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrInvalidCredentials) {
			status = http.StatusUnauthorized
		}

		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
