package repository

import (
	"database/sql"
	"errors"
	"strings"
	"fmt"
	"log"
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
func (r *UserRepository) GetUsersRepo(search, sortBy, order string, limit, offset int) ([]models.User, error) {
	query := fmt.Sprintf(`
		SELECT id, username, email, created_at
		FROM users
		WHERE username ILIKE $1 OR email ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) CountUsersRepo(search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE username ILIKE $1 OR email ILIKE $1`
	err := r.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

