package repository

import (
	"database/sql"
	"errors"
	"strings"
	"prak4/app/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) GetByUsernameOrEmail(identifier string) (*models.User, string, error) {
	u := models.User{}
	var hash string
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, role
		FROM users
		WHERE username = $1 OR email = $1
	`, identifier).Scan(&u.ID, &u.Username, &u.Email, &hash, &u.Role)
	if err != nil {
		return nil, "", err
	}
	return &u, hash, nil
}

func (r *UserRepository) ExistsByUsernameOrEmail(username, email string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM users WHERE username = $1 OR email = $2
		)
	`, username, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) Create(username, email, passwordHash, role string) (*models.User, error) {
	role = strings.ToLower(role)
	if role != "admin" && role != "user" {
		return nil, errors.New("role tidak valid")
	}
	var u models.User
	err := r.DB.QueryRow(`
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, role
	`, username, email, passwordHash, role).Scan(&u.ID, &u.Username, &u.Email, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
