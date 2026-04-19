package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type TweetService interface {
	CreateTweet(ctx context.Context, userID int32, content string) (*db.Tweet, error)
	GetUserTweets(ctx context.Context, userID int32) ([]db.Tweet, error)
	GetAllTweets(ctx context.Context, limit int32) ([]db.Tweet, error)
	DeleteTweet(ctx context.Context, tweetID int32, userID int32) error
}

type tweetService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewTweetService(database *sql.DB, queries *db.Queries) TweetService {
	return &tweetService{
		db:      database,
		queries: queries,
	}
}

// CreateTweet はツイートを作成
func (s *tweetService) CreateTweet(ctx context.Context, userID int32, content string) (*db.Tweet, error) {
	// コンテンツの検証
	if content == "" {
		return nil, &ValidationError{Message: "ツイート内容を入力してください"}
	}
	
	if len([]rune(content)) > 140 {
		return nil, &ValidationError{Message: "ツイートは140文字以内で入力してください"}
	}

	// ツイートを作成
	tweet, err := s.queries.CreateTweet(ctx, db.CreateTweetParams{
		UserID:  userID,
		Content: content,
	})
	if err != nil {
		return nil, &ServiceError{Message: "ツイートの投稿に失敗しました"}
	}

	return &tweet, nil
}

// GetUserTweets は特定ユーザーのツイート一覧を取得
func (s *tweetService) GetUserTweets(ctx context.Context, userID int32) ([]db.Tweet, error) {
	tweets, err := s.queries.GetTweetsByUserID(ctx, userID)
	if err != nil {
		return nil, &ServiceError{Message: "ツイートの取得に失敗しました"}
	}
	return tweets, nil
}

// GetAllTweets は全ツイートを取得
func (s *tweetService) GetAllTweets(ctx context.Context, limit int32) ([]db.Tweet, error) {
	if limit <= 0 {
		limit = 20 // デフォルト20件
	}
	
	tweets, err := s.queries.GetAllTweets(ctx, limit)
	if err != nil {
		return nil, &ServiceError{Message: "ツイートの取得に失敗しました"}
	}
	return tweets, nil
}

// DeleteTweet はツイートを削除（自分のツイートのみ）
func (s *tweetService) DeleteTweet(ctx context.Context, tweetID int32, userID int32) error {
	err := s.queries.DeleteTweet(ctx, db.DeleteTweetParams{
		ID:     tweetID,
		UserID: userID,
	})
	if err != nil {
		return &ServiceError{Message: "ツイートの削除に失敗しました"}
	}
	return nil
}
