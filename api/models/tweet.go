package models

import (
	"time"
)

type Tweet struct {
	ID           int64     `json:"id" db:"id"`
	AuthorID     int64     `json:"author_id" db:"author_id"`
	Author       *User     `json:"author,omitempty" db:"-"`
	Text         string    `json:"text" db:"text"`
	ReplyToID    *int64    `json:"reply_to_id,omitempty" db:"reply_to_id"`
	RetweetOfID  *int64    `json:"retweet_of_id,omitempty" db:"retweet_of_id"`
	LikeCount    int       `json:"like_count" db:"like_count"`
	RetweetCount int       `json:"retweet_count" db:"retweet_count"`
	ReplyCount   int       `json:"reply_count" db:"reply_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type CreateTweetRequest struct {
	Text        string `json:"text" binding:"required,max=280"`
	ReplyToID   *int64 `json:"reply_to_id,omitempty"`
	RetweetOfID *int64 `json:"retweet_of_id,omitempty"`
}

type TimelineResponse struct {
	Tweets     []*Tweet `json:"tweets"`
	NextCursor *int64   `json:"next_cursor,omitempty"`
	HasMore    bool     `json:"has_more"`
}