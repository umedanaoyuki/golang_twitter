package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type BookmarkService interface {
	CreateBookmark(ctx context.Context, userID int32, tweetID int32) error
	DeleteBookmark(ctx context.Context, userID int32, tweetID int32) error
	GetBookmarksByUserId(ctx context.Context, userID int32) ([]db.GetBookmarksByUserIdRow, error)
}

type bookmarkService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewBookmarkService(db *sql.DB, queries *db.Queries) BookmarkService {
	return &bookmarkService{
		db:      db,
		queries: queries,
	}
}

func (s *bookmarkService) CreateBookmark(ctx context.Context, userID int32, tweetID int32) error {
	_, err := s.queries.CreateBookmark(ctx, db.CreateBookmarkParams{
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

func (s *bookmarkService) DeleteBookmark(ctx context.Context, userID int32, tweetID int32) error {
	return s.queries.DeleteBookmark(ctx, db.DeleteBookmarkParams{
		UserID:  userID,
		TweetID: tweetID,
	})
}

func (s *bookmarkService) GetBookmarksByUserId(ctx context.Context, userID int32) ([]db.GetBookmarksByUserIdRow, error) {
	bookmarks, err := s.queries.GetBookmarksByUserId(ctx, userID)
	if err != nil {
		return nil, &ServiceError{Message: "ブックマークの取得に失敗しました"}
	}
	return bookmarks, nil
}