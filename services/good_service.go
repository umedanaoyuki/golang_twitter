package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type GoodService interface {
	CreateGood(ctx context.Context, userID int32, tweetID int32) error
	DeleteGood(ctx context.Context, userID int32, tweetID int32) error
}

type goodService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewGoodService(db *sql.DB, queries *db.Queries) GoodService {
	return &goodService{
		db:      db,
		queries: queries,
	}
}

func (s *goodService) CreateGood(ctx context.Context, userID int32, tweetID int32) error {
	_, err := s.queries.CreateLike(ctx, db.CreateLikeParams{
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

func (s *goodService) DeleteGood(ctx context.Context, userID int32, tweetID int32) error {
	return s.queries.DeleteLike(ctx, db.DeleteLikeParams{
		UserID:  userID,
		TweetID: tweetID,
	})
}