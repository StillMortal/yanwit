package repository

import (
	"database/sql"
	"fmt"
)

// Follow подписывает пользователя на другого
func Follow(followerID, followeeID int64) error {
	query := `INSERT INTO follows (follower_id, followee_id) VALUES ($1, $2)`
	_, err := DB.Exec(query, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("failed to follow: %w", err)
	}
	return nil
}

// Unfollow отписывает пользователя
func Unfollow(followerID, followeeID int64) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2`
	_, err := DB.Exec(query, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("failed to unfollow: %w", err)
	}
	return nil
}

// IsFollowing проверяет, подписан ли пользователь
func IsFollowing(followerID, followeeID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND followee_id = $2)`
	err := DB.QueryRow(query, followerID, followeeID).Scan(&exists)
	return exists, err
}

// GetFollowers возвращает список подписчиков пользователя
func GetFollowers(userID int64) ([]int64, error) {
	var followers []int64
	query := `SELECT follower_id FROM follows WHERE followee_id = $1`
	err := DB.Select(&followers, query, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	return followers, nil
}

// GetFollowees возвращает список пользователей, на кого подписан
func GetFollowees(userID int64) ([]int64, error) {
	var followees []int64
	query := `SELECT followee_id FROM follows WHERE follower_id = $1`
	err := DB.Select(&followees, query, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get followees: %w", err)
	}
	return followees, nil
}