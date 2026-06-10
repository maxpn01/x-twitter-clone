package user

import (
	"errors"
	"net/mail"
	"unicode"
)

func ValidateSignupInput(input SignupInput) error {
	if err := validateEmail(input.Email); err != nil {
		return err
	}

	if err := validateUsername(input.Username); err != nil {
		return err
	}

	if err := validateFullname(input.Fullname); err != nil {
		return err
	}

	if err := validatePassword(input.Password); err != nil {
		return err
	}

	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	if len(email) > 255 {
		return errors.New("email must be at most 255 characters")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email is invalid")
	}

	return nil
}

func validateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if len(username) > 20 {
		return errors.New("username must be at most 20 characters")
	}

	for _, r := range username {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '_' {
			continue
		}

		return errors.New("username may only contain letters, numbers, and underscores")
	}

	return nil
}

func validateFullname(fullname string) error {
	if len(fullname) < 3 {
		return errors.New("fullname must be at least 3 characters")
	}

	if len(fullname) > 50 {
		return errors.New("fullname must be at most 50 characters")
	}

	for _, r := range fullname {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r == ' ' {
			continue
		}

		return errors.New("fullname may only contain letters and spaces")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if len(password) > 72 {
		return errors.New("password must be at most 72 characters")
	}

	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSymbol bool

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSymbol {
		return errors.New("password must contain at least one uppercase letter, lowercase letter, number, and symbol")
	}

	return nil
}
