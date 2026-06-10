package repository

import (
	"fmt"

	"yanwit/api/models"
)

func CreateTweet(tweet *models.Tweet) error {
	query := `
		INSERT INTO tweets (author_id, text, reply_to_id, retweet_of_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := DB.QueryRow(query, tweet.AuthorID, tweet.Text, tweet.ReplyToID, tweet.RetweetOfID).Scan(
		&tweet.ID, &tweet.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create tweet: %w", err)
	}
	return nil
}

func GetTweetByID(id int64) (*models.Tweet, error) {
	var tweet models.Tweet
	query := `SELECT id, author_id, text, reply_to_id, retweet_of_id,
	                 like_count, retweet_count, reply_count, created_at
	          FROM tweets WHERE id = $1`
	err := DB.Get(&tweet, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweet: %w", err)
	}
	return &tweet, nil
}

// GetTweetsByAuthorIDs возвращает твиты авторов из списка (для ленты)
func GetTweetsByAuthorIDs(authorIDs []int64, limit, offset int) ([]*models.Tweet, error) {
	if len(authorIDs) == 0 {
		return []*models.Tweet{}, nil
	}

	query := `
		SELECT id, author_id, text, reply_to_id, retweet_of_id,
		       like_count, retweet_count, reply_count, created_at
		FROM tweets
		WHERE author_id = ANY($1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var tweets []*models.Tweet
	err := DB.Select(&tweets, query, authorIDs, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweets: %w", err)
	}
	return tweets, nil
}

// GetTweetsByUserID возвращает твиты конкретного пользователя
func GetTweetsByUserID(userID int64, limit, offset int) ([]*models.Tweet, error) {
	var tweets []*models.Tweet
	query := `
		SELECT id, author_id, text, reply_to_id, retweet_of_id,
		       like_count, retweet_count, reply_count, created_at
		FROM tweets
		WHERE author_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := DB.Select(&tweets, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tweets: %w", err)
	}
	return tweets, nil
}

// GetTweetsByIDs возвращает твиты по списку ID
func GetTweetsByIDs(ids []int64) ([]*models.Tweet, error) {
	if len(ids) == 0 {
		return []*models.Tweet{}, nil
	}

	query := `
		SELECT id, author_id, text, reply_to_id, retweet_of_id,
		       like_count, retweet_count, reply_count, created_at
		FROM tweets
		WHERE id = ANY($1)
		ORDER BY created_at DESC
	`

	var tweets []*models.Tweet
	err := DB.Select(&tweets, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweets by IDs: %w", err)
	}
	return tweets, nil
}