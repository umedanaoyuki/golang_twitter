package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID int32, tweetID int32) error
	DeleteComment(ctx context.Context, userID int32, tweetID int32) error
	GetCommentsByTweetIDWithCursor(ctx context.Context, tweetID int32, cursor *int32, limit int32) ([]db.Comment, error)
}

type commentService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewCommentService(db *sql.DB, queries *db.Queries) CommentService {
	return &commentService{
		db:      db,
		queries: queries,
	}
}

func (s *commentService) CreateComment(ctx context.Context, userID int32, tweetID int32) error {
	_, err := s.queries.CreateComment(ctx, db.CreateCommentParams{
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

func (s *commentService) DeleteComment(ctx context.Context, userID int32, tweetID int32) error {
	return s.queries.DeleteComment(ctx, db.DeleteCommentParams{
		UserID:  userID,
		TweetID: tweetID,
	})
}

func (s *commentService) GetCommentsByTweetIDWithCursor(ctx context.Context, tweetID int32, cursor *int32, limit int32) ([]db.Comment, error) {
	cursorValue := int32(0)
	if cursor != nil {
		cursorValue = *cursor
	}

	comments, err := s.queries.GetCommentsByTweetIDWithCursor(ctx, db.GetCommentsByTweetIDWithCursorParams{
		TweetID: tweetID,
		Column2: cursorValue,
		Limit:   limit,
	})
	if err != nil {
		return nil, &ServiceError{Message: "コメントの取得に失敗しました"}
	}

	return comments, nil
}