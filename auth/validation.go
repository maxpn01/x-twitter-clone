package auth

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/maxpn01/x-twitter-clone/apperror"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func normalizeAuthInput(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func validateEmail(email string) error {
	if email == "" {
		return validationError("email is required")
	}
	if len(email) > 255 {
		return validationError("email must be 255 characters or fewer")
	}
	if strings.ContainsAny(email, " \t\r\n") {
		return validationError("email must not contain whitespace")
	}

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return validationError("email must be a valid email address")
	}

	parts := strings.Split(address.Address, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" || !strings.Contains(parts[1], ".") {
		return validationError("email must be a valid email address")
	}

	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return validationError("username is required")
	}
	if len(username) < 3 {
		return validationError("username must be at least 3 characters")
	}
	if len(username) > 20 {
		return validationError("username must be 20 characters or fewer")
	}
	if !usernamePattern.MatchString(username) {
		return validationError("username can only contain letters, numbers and underscores")
	}
	if strings.HasPrefix(username, "_") || strings.HasSuffix(username, "_") {
		return validationError("username must not start or end with an underscore")
	}
	if strings.Contains(username, "__") {
		return validationError("username must not contain consecutive underscores")
	}

	return nil
}

func validatePassword(password, email, username string) error {
	if password == "" {
		return validationError("password is required")
	}
	if len(password) < 8 {
		return validationError("password must be at least 8 characters")
	}
	if len(password) > 72 {
		return validationError("password must be at max 72 characters")
	}

	var hasLower, hasUpper, hasDigit, hasSymbol bool
	for _, r := range password {
		if unicode.IsSpace(r) {
			return validationError("password must not contain whitespace")
		}
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}

	if !hasLower || !hasUpper || !hasDigit || !hasSymbol {
		return validationError("password must include uppercase, lowercase, number and symbol characters")
	}

	lowerPassword := strings.ToLower(password)
	if username != "" && strings.Contains(lowerPassword, strings.ToLower(username)) {
		return validationError("password must not contain the username")
	}

	emailLocalPart := strings.Split(email, "@")[0]
	if utf8.RuneCountInString(emailLocalPart) >= 3 && strings.Contains(lowerPassword, strings.ToLower(emailLocalPart)) {
		return validationError("password must not contain the email name")
	}

	return nil
}

func validationError(message string) error {
	return apperror.New(apperror.CodeValidation, message)
}
