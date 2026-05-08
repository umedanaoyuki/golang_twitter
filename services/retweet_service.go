package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type RetweetService interface {
	CreateRetweet(ctx context.Context, userID int32, tweetID int32) error
	DeleteRetweet(ctx context.Context, userID int32, tweetID int32) error
	GetUserRetweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Retweet, error)
}

type retweetService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewRetweetService(db *sql.DB, queries *db.Queries) RetweetService {
	return &retweetService{
		db:      db,
		queries: queries,
	}
}

func (s *retweetService) CreateRetweet(ctx context.Context, userID int32, tweetID int32) error {
	_, err := s.queries.CreateRetweet(ctx, db.CreateRetweetParams{
		UserID:  userID,
		TweetID: tweetID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func (s *retweetService) DeleteRetweet(ctx context.Context, userID int32, tweetID int32) error {
	return s.queries.DeleteRetweet(ctx, db.DeleteRetweetParams{
		UserID:  userID,
		TweetID: tweetID,
	})
}

// ユーザーのリツイート一覧を取得（ページネーションつき）
func (s *retweetService) GetUserRetweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Retweet, error) {
	cursorValue := int32(0)
	if cursor != nil {
		cursorValue = *cursor
	}
	return s.queries.GetUserRetweetsWithCursor(ctx, db.GetUserRetweetsWithCursorParams{
		UserID:  userID,
		Column2: cursorValue,
		Limit:   limit,
	})
}