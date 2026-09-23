package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"
)

type FollowService interface {
	CreateFollow(ctx context.Context, userID int32, followedUserID int32) error
	DeleteFollow(ctx context.Context, userID int32, followedUserID int32) error
	GetFollowersByUserIdWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Follow, error)
	GetFollowingByUserIdWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Follow, error)
}

type followService struct {
	db                  *sql.DB
	queries             *db.Queries
	notificationService NotificationService
}

func NewFollowService(db *sql.DB, queries *db.Queries, notificationService NotificationService) FollowService {
	return &followService{
		db:                  db,
		queries:             queries,
		notificationService: notificationService,
	}
}

func (s *followService) CreateFollow(ctx context.Context, userID int32, followedUserID int32) error {
	if userID == followedUserID {
		return &ValidationError{Message: "自分自身をフォローすることはできません"}
	}

	_, err := s.queries.CreateFollow(ctx, db.CreateFollowParams{
		UserID:         userID,
		FollowedUserID: followedUserID,
	})
	if err != nil {
		// ON CONFLICT DO NOTHING により既にフォロー済みの場合は ErrNoRows になる（通知も作らない）
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	// 新規にフォローされた場合のみフォローされたユーザーへ通知
	s.notificationService.NotifyFollow(ctx, userID, followedUserID)
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