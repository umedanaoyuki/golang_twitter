package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type FollowService interface {
	CreateFollow(ctx context.Context, userID int32, followedUserID int32) error
	DeleteFollow(ctx context.Context, userID int32, followedUserID int32) error
	GetFollowersByUserIdWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Follow, error)
	GetFollowingByUserIdWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Follow, error)
}

type followService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewFollowService(db *sql.DB, queries *db.Queries) FollowService {
	return &followService{
		db:      db,
		queries: queries,
	}
}

func (s *followService) CreateFollow(ctx context.Context, userID int32, followedUserID int32) error {
	_, err := s.queries.CreateFollow(ctx, db.CreateFollowParams{
		UserID: userID,
		FollowedUserID: followedUserID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *followService) DeleteFollow(ctx context.Context, userID int32, followedUserID int32) error {
	return s.queries.DeleteFollow(ctx, db.DeleteFollowParams{
		UserID:         userID,
		FollowedUserID: followedUserID,
	})
}

func (s *followService) GetFollowersByUserIdWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Follow, error) {
	cursorValue := int32(0)
	if cursor != nil {
		cursorValue = *cursor
	}
	return s.queries.GetFollowersByUserIdWithCursor(ctx, db.GetFollowersByUserIdWithCursorParams{
		FollowedUserID: userID,
		Column2:        cursorValue,
		Limit:          limit,
	})
}

func (s *followService) GetFollowingByUserIdWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Follow, error) {
	cursorValue := int32(0)
	if cursor != nil {
		cursorValue = *cursor
	}
	return s.queries.GetFollowingByUserIdWithCursor(ctx, db.GetFollowingByUserIdWithCursorParams{
		UserID:  userID,
		Column2: cursorValue,
		Limit:   limit,
	})
}