package apperror

type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

const (
	CodeBadRequest            = "BAD_REQUEST"
	CodeValidation            = "VALIDATION_ERROR"
	CodeEmailAlreadyExists    = "EMAIL_ALREADY_EXISTS"
	CodeUsernameAlreadyExists = "USERNAME_ALREADY_EXISTS"
	CodeNotFound              = "NOT_FOUND"
	CodeInternal              = "INTERNAL_ERROR"
)
