package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"
	"golang_twitter/infrastructure/storage"
)

var ErrTweetNotFound = errors.New("Tweetが見つかりませんでした")

type TweetDetail struct {
	db.Tweet
	// いいね数
	LikeCount int64 `json:"like_count"`
	// リツイート数
	RetweetCount int64 `json:"retweet_count"`
	// コメント数
	CommentCount int64 `json:"comment_count"`
}

type TweetService interface {
	CreateTweet(ctx context.Context, userID int32, content string) (*db.Tweet, error)
	DeleteTweet(ctx context.Context, userID int32, id int32) error
	PresignImageUpload(ctx context.Context, userID int32, contentType string, size int64) (*storage.PresignUploadResult, error)
	CompleteImageUpload(ctx context.Context, userID int32, key string) (*db.Tweet, error)
	GetTweetByID(ctx context.Context, id int32) (*TweetDetail, error)
	GetUserTweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]TweetDetail, error)
	// GetTweetDetailsByIDs は指定 ID のツイートを一括取得し、いいね数・リツイート数・コメント数を付与する
	GetTweetDetailsByIDs(ctx context.Context, ids []int32) (map[int32]TweetDetail, error)
}

type tweetService struct {
	db      *sql.DB
	queries *db.Queries
	imageStorage storage.ImageStorage
}

func NewTweetService(database *sql.DB, queries *db.Queries, imageStorage storage.ImageStorage) TweetService {
	return &tweetService{
		db:      database,
		queries: queries,
		imageStorage: imageStorage,
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
		UserID:   userID,
		Content:  content,
		ImageUrl: "",
	})
	if err != nil {
		return nil, &ServiceError{Message: "ツイートの投稿に失敗しました"}
	}

	return &tweet, nil
}

func (s *tweetService) DeleteTweet(ctx context.Context, userID int32, id int32) error {
	err := s.queries.DeleteTweet(ctx, db.DeleteTweetParams{
		ID: id,
		UserID: userID,
	})
	if err != nil {
		return &ServiceError{Message: "ツイートの削除に失敗しました"}
	}

	return nil
}

func (s *tweetService) PresignImageUpload(ctx context.Context, userID int32, contentType string, size int64) (*storage.PresignUploadResult, error) {
	if s.imageStorage == nil {
		return nil, &ServiceError{Message: "画像ストレージが設定されていません"}
	}

	result, err := s.imageStorage.PresignUpload(userID, contentType, size)
	if err != nil {
		return nil, &ValidationError{Message: err.Error()}
	}

	return result, nil
}

func (s *tweetService) CompleteImageUpload(ctx context.Context, userID int32, key string) (*db.Tweet, error) {
	if s.imageStorage == nil {
		return nil, &ServiceError{Message: "画像ストレージが設定されていません"}
	}

	if err := storage.ValidateKeyForUser(userID, key); err != nil {
		return nil, &ValidationError{Message: err.Error()}
	}

	imageURL, err := s.imageStorage.ConfirmUpload(key)
	if err != nil {
		return nil, &ValidationError{Message: err.Error()}
	}

	tweet, err := s.queries.CreateTweet(ctx, db.CreateTweetParams{
		UserID:   userID,
		Content:  "",
		ImageUrl: imageURL,
	})
	if err != nil {
		return nil, &ServiceError{Message: "画像の投稿に失敗しました"}
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

	tweetIDs := make([]int32, 0, len(tweets))
	for _, t := range tweets {
		tweetIDs = append(tweetIDs, t.ID)
	}

	// コメント数はツイート単位で数えると件数分クエリが増えるため、一括で取得する
	commentRows, err := s.queries.CountCommentsByTweetIDs(ctx, tweetIDs)
	if err != nil {
		return nil, &ServiceError{Message: "コメント数の取得に失敗しました"}
	}
	commentByTweetID := make(map[int32]int64, len(commentRows))
	for _, row := range commentRows {
		commentByTweetID[row.TweetID] = row.CommentCount
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
		details = append(details, TweetDetail{
			Tweet:        t,
			LikeCount:    n,
			RetweetCount: rt,
			CommentCount: commentByTweetID[t.ID],
		})
	}
	return details, nil
}

func (s *tweetService) GetTweetByID(ctx context.Context, id int32) (*TweetDetail, error) {
	tweet, err := s.queries.GetTweetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTweetNotFound
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

	commentCount, err := s.queries.CountCommentsByTweetID(ctx, id)
	if err != nil {
		return nil, &ServiceError{Message: "コメント数の取得に失敗しました"}
	}

	return &TweetDetail{
		Tweet:        tweet,
		LikeCount:    likeCount,
		RetweetCount: retweetCount,
		CommentCount: commentCount,
	}, nil
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

	likeRows, err := s.queries.CountLikesByTweetIDs(ctx, unique)
	if err != nil {
		return nil, &ServiceError{Message: "いいね数の取得に失敗しました"}
	}
	retweetRows, err := s.queries.CountRetweetsByTweetIDs(ctx, unique)
	if err != nil {
		return nil, &ServiceError{Message: "リツイート数の取得に失敗しました"}
	}
	commentRows, err := s.queries.CountCommentsByTweetIDs(ctx, unique)
	if err != nil {
		return nil, &ServiceError{Message: "コメント数の取得に失敗しました"}
	}

	likeByTweetID := make(map[int32]int64, len(likeRows))
	for _, row := range likeRows {
		likeByTweetID[row.TweetID] = row.LikeCount
	}
	retweetByTweetID := make(map[int32]int64, len(retweetRows))
	for _, row := range retweetRows {
		retweetByTweetID[row.TweetID] = row.RetweetCount
	}
	commentByTweetID := make(map[int32]int64, len(commentRows))
	for _, row := range commentRows {
		commentByTweetID[row.TweetID] = row.CommentCount
	}

	out := make(map[int32]TweetDetail, len(tweets))
	for _, tweet := range tweets {
		out[tweet.ID] = TweetDetail{
			Tweet:        tweet,
			LikeCount:    likeByTweetID[tweet.ID],
			RetweetCount: retweetByTweetID[tweet.ID],
			CommentCount: commentByTweetID[tweet.ID],
		}
	}
	return out, nil
}