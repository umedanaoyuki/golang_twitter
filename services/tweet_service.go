package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type TweetService interface {
	CreateTweet(ctx context.Context, userID int32, content string) (*db.Tweet, error)
	GetTweetByID(ctx context.Context, id int32) (*db.Tweet, error)
	GetUserTweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Tweet, error)
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

// userIDのツイート一覧を取得
func (s *tweetService) GetUserTweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]db.Tweet, error) {
	// cursorが指定されていない場合は0を使用（SQLで全件取得）
	cursorValue := int32(0)
	if cursor != nil {
		cursorValue = *cursor
	}

	// カーソルベースでツイートを取得
	tweets, err := s.queries.GetTweetsByUserIDWithCursor(ctx, db.GetTweetsByUserIDWithCursorParams{
		UserID:  userID,
		Column2: cursorValue,
		Limit:   limit,
	})
	if err != nil {
		return nil, &ServiceError{Message: "ツイートの取得に失敗しました"}
	}

	return tweets, nil
}

func (s *tweetService) GetTweetByID(ctx context.Context, id int32) (*db.Tweet, error) {
	tweet, err := s.queries.GetTweetByID(ctx, id)
	if err != nil {
		return nil, &ServiceError{Message: "Tweetが見つかりませんでした"}
	}

	return &tweet, nil
}