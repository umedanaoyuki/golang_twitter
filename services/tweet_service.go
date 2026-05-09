package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"
)

type TweetDetail struct {
	db.Tweet
	// いいね数
	LikeCount int64 `json:"like_count"`
	// リツイート数
	RetweetCount int64 `json:"retweet_count"`
}

type TweetService interface {
	CreateTweet(ctx context.Context, userID int32, content string) (*db.Tweet, error)
	GetTweetByID(ctx context.Context, id int32) (*TweetDetail, error)
	GetUserTweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]TweetDetail, error)
	// GetTweetDetailsByIDs は指定 ID のツイートを一括取得し、いいね数・リツイート数を付与する
	GetTweetDetailsByIDs(ctx context.Context, ids []int32) (map[int32]TweetDetail, error)
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
func (s *tweetService) GetUserTweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]TweetDetail, error) {
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

	if len(tweets) == 0 {
		return []TweetDetail{}, nil
	}

	details := make([]TweetDetail, 0, len(tweets))
	for _, t := range tweets {
		n, err := s.queries.CountLikesByTweetID(ctx, t.ID)
		if err != nil {
			return nil, &ServiceError{Message: "いいね数の取得に失敗しました"}
		}
		rt, err := s.queries.CountRetweetsByTweetID(ctx, t.ID)
		if err != nil {
			return nil, &ServiceError{Message: "リツイート数の取得に失敗しました"}
		}
		details = append(details, TweetDetail{Tweet: t, LikeCount: n, RetweetCount: rt})
	}
	return details, nil
}

func (s *tweetService) GetTweetByID(ctx context.Context, id int32) (*TweetDetail, error) {
	tweet, err := s.queries.GetTweetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Tweetが見つかりませんでした");
		}
		return nil, &ServiceError{Message: "Tweetの取得に失敗しました"}
	}

	likeCount, err := s.queries.CountLikesByTweetID(ctx, id)
	if err != nil {
		return nil, &ServiceError{Message: "いいね数の取得に失敗しました"}
	}

	retweetCount, err := s.queries.CountRetweetsByTweetID(ctx, id)
	if err != nil {
		return nil, &ServiceError{Message: "リツイート数の取得に失敗しました"}
	}
	return &TweetDetail{Tweet: tweet, LikeCount: likeCount, RetweetCount: retweetCount}, nil
}

func (s *tweetService) GetTweetDetailsByIDs(ctx context.Context, ids []int32) (map[int32]TweetDetail, error) {
	if len(ids) == 0 {
		return map[int32]TweetDetail{}, nil
	}
	seen := make(map[int32]struct{}, len(ids))
	unique := make([]int32, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	tweets, err := s.queries.GetTweetsByIDs(ctx, unique)
	if err != nil {
		return nil, &ServiceError{Message: "ツイートの取得に失敗しました"}
	}

	out := make(map[int32]TweetDetail, len(tweets))
	for _, t := range tweets {
		n, err := s.queries.CountLikesByTweetID(ctx, t.ID)
		if err != nil {
			return nil, &ServiceError{Message: "いいね数の取得に失敗しました"}
		}
		rt, err := s.queries.CountRetweetsByTweetID(ctx, t.ID)
		if err != nil {
			return nil, &ServiceError{Message: "リツイート数の取得に失敗しました"}
		}
		out[t.ID] = TweetDetail{Tweet: t, LikeCount: n, RetweetCount: rt}
	}
	return out, nil
}