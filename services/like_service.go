package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type LikeService interface {
	CreateLike(ctx context.Context, userID int32, tweetID int32) error
	DeleteLike(ctx context.Context, userID int32, tweetID int32) error
}

type likeService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewLikeService(db *sql.DB, queries *db.Queries) LikeService {
	return &likeService{
		db:      db,
		queries: queries,
	}
}

func (s *likeService) CreateLike(ctx context.Context, userID int32, tweetID int32) error {
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

func (s *likeService) DeleteLike(ctx context.Context, userID int32, tweetID int32) error {
	return s.queries.DeleteLike(ctx, db.DeleteLikeParams{
		UserID:  userID,
		TweetID: tweetID,
	})
}