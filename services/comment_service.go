package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID int32, tweetID int32, content string) (*db.Comment, error)
	DeleteComment(ctx context.Context, userID int32, tweetID int32, commentID int32) error
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

func (s *commentService) CreateComment(ctx context.Context, userID int32, tweetID int32, content string) (*db.Comment, error) {
	if content == "" {
		return nil, &ValidationError{Message: "コメント内容を入力してください"}
	}
	if len([]rune(content)) > 140 {
		return nil, &ValidationError{Message: "コメントは140文字以内で入力してください"}
	}

	comment, err := s.queries.CreateComment(ctx, db.CreateCommentParams{
		UserID:  userID,
		TweetID: tweetID,
		Content: content,
	})
	if err != nil {
		return nil, &ServiceError{Message: "コメントの作成に失敗しました"}
	}
	return &comment, nil
}

func (s *commentService) DeleteComment(ctx context.Context, userID int32, tweetID int32, commentID int32) error {
	_, err := s.queries.DeleteComment(ctx, db.DeleteCommentParams{
		ID:      commentID,
		UserID:  userID,
		TweetID: tweetID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &ValidationError{Message: "コメントが見つかりません"}
		}
		return &ServiceError{Message: "コメントの削除に失敗しました"}
	}
	return nil
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
