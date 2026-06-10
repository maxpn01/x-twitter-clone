package user

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/maxpn01/x-twitter-clone/models"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid email, username or password")
)

type CreateUserInput struct {
	Email        string
	Username     string
	Fullname     string
	PasswordHash string
}

type UserRepository interface {
	CreateUser(input CreateUserInput) (models.User, error)
	GetUserByID(id string) (models.User, error)
	GetUserByEmailOrUsername(emailOrUsername string) (models.User, error)
}

type PostgresUserRepository struct {
	DB *sql.DB
}

func (r *PostgresUserRepository) CreateUser(input CreateUserInput) (models.User, error) {
	query := `
		INSERT INTO users (email, username, fullname, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, username, fullname, password_hash, created_at, updated_at
	`

	var user models.User

	err := r.DB.QueryRow(
		query,
		input.Email,
		input.Username,
		input.Fullname,
		input.PasswordHash,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Fullname,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, translateCreateUserError(err)
	}

	return user, nil
}

func translateCreateUserError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Code != pgerrcode.UniqueViolation {
		return err
	}

	switch pgErr.ConstraintName {
	case "users_email_key":
		return ErrEmailAlreadyExists
	case "users_username_key":
		return ErrUsernameAlreadyExists
	default:
		return err
	}
}

func (r *PostgresUserRepository) GetUserByID(id string) (models.User, error) {
	return models.User{}, nil
}

func (r *PostgresUserRepository) GetUserByEmailOrUsername(emailOrUsername string) (models.User, error) {
	query := `SELECT id, email, username, fullname, password_hash, created_at, updated_at
			  FROM users
			  WHERE email = $1 OR username = $1`

	var user models.User

	err := r.DB.QueryRow(
		query,
		emailOrUsername,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Fullname,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, translateGetUserError(err)
	}

	return user, nil
}

func translateGetUserError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	return err
}
