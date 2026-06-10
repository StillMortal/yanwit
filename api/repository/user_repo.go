package repository

import (
	"database/sql"
	"fmt"

	"yanwit/api/models"
)

func CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := DB.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, avatar_url, bio, created_at, updated_at
	          FROM users WHERE username = $1`
	err := DB.Get(&user, query, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func GetUserByID(id int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, avatar_url, bio, created_at, updated_at
	          FROM users WHERE id = $1`
	err := DB.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}