package services

import (
	"database/sql"
	"time"

	db "golang_twitter/db/sqlc"
)

// Tweet はAPIレスポンス用のツイート表現
type Tweet struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  *string   `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TweetDetail struct {
	Tweet
	LikeCount    int64 `json:"like_count"`
	RetweetCount int64 `json:"retweet_count"`
}

func imageURLFromNull(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

func tweetFromCreateRow(r db.CreateTweetRow) Tweet {
	return Tweet{
		ID:        r.ID,
		UserID:    r.UserID,
		Content:   r.Content,
		ImageURL:  imageURLFromNull(r.ImageUrl),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func tweetFromCreateWithImageRow(r db.CreateTweetWithImageRow) Tweet {
	return Tweet{
		ID:        r.ID,
		UserID:    r.UserID,
		Content:   r.Content,
		ImageURL:  imageURLFromNull(r.ImageUrl),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func tweetFromGetByIDRow(r db.GetTweetByIDRow) Tweet {
	return Tweet{
		ID:        r.ID,
		UserID:    r.UserID,
		Content:   r.Content,
		ImageURL:  imageURLFromNull(r.ImageUrl),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func tweetFromCursorRow(r db.GetTweetsByUserIDWithCursorRow) Tweet {
	return Tweet{
		ID:        r.ID,
		UserID:    r.UserID,
		Content:   r.Content,
		ImageURL:  imageURLFromNull(r.ImageUrl),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func tweetFromIDsRow(r db.GetTweetsByIDsRow) Tweet {
	return Tweet{
		ID:        r.ID,
		UserID:    r.UserID,
		Content:   r.Content,
		ImageURL:  imageURLFromNull(r.ImageUrl),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
