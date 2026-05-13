package repository

import (
	"database/sql"

	"github.com/maxpn01/x-twitter-clone/models"
)

type UserRepository interface {
	CreateUser(user models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user models.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, username, fullname, password_hash) VALUES ($1, $2, $3, $4)", user.Email, user.Username, user.Fullname, user.PasswordHash)

	return err
}

func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	row := r.db.QueryRow("SELECT id, email, username, fullname, password_hash, created_at, updated_at FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Fullname, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User

	row := r.db.QueryRow("SELECT id, email, username, fullname, password_hash, created_at, updated_at FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Fullname, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
