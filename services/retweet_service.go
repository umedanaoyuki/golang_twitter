package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type RetweetService interface {
	CreateRetweet(ctx context.Context, userID int32, tweetID int32) error
	DeleteRetweet(ctx context.Context, userID int32, tweetID int32) error
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