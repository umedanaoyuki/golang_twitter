package services

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"path/filepath"
	"strings"

	db "golang_twitter/db/sqlc"
	"golang_twitter/infrastructure/storage"
)

type TweetService interface {
	CreateTweet(ctx context.Context, userID int32, content string) (*Tweet, error)
	CreateTweetWithImage(ctx context.Context, userID int32, content string, filename string, contentType string, r io.Reader, size int64) (*Tweet, error)
	GetTweetByID(ctx context.Context, id int32) (*TweetDetail, error)
	GetUserTweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]TweetDetail, error)
	GetTweetDetailsByIDs(ctx context.Context, ids []int32) (map[int32]TweetDetail, error)
}

type tweetService struct {
	db           *sql.DB
	queries      *db.Queries
	imageStorage storage.ImageStorage
}

func NewTweetService(database *sql.DB, queries *db.Queries, imageStorage storage.ImageStorage) TweetService {
	return &tweetService{
		db:           database,
		queries:      queries,
		imageStorage: imageStorage,
	}
}

func (s *tweetService) validateContent(content string, allowEmpty bool) error {
	if !allowEmpty && content == "" {
		return &ValidationError{Message: "ツイート内容を入力してください"}
	}
	if len([]rune(content)) > 140 {
		return &ValidationError{Message: "ツイートは140文字以内で入力してください"}
	}
	return nil
}

// CreateTweet はテキストのみのツイートを作成
func (s *tweetService) CreateTweet(ctx context.Context, userID int32, content string) (*Tweet, error) {
	if err := s.validateContent(content, false); err != nil {
		return nil, err
	}

	row, err := s.queries.CreateTweet(ctx, db.CreateTweetParams{
		UserID:  userID,
		Content: content,
	})
	if err != nil {
		return nil, &ServiceError{Message: "ツイートの投稿に失敗しました"}
	}

	tweet := tweetFromCreateRow(row)
	return &tweet, nil
}

// CreateTweetWithImage は画像付きツイートを作成（本文は任意）
func (s *tweetService) CreateTweetWithImage(
	ctx context.Context,
	userID int32,
	content string,
	filename string,
	contentType string,
	r io.Reader,
	size int64,
) (*Tweet, error) {
	if err := s.validateContent(content, true); err != nil {
		return nil, err
	}

	resolvedType, err := resolveImageContentType(contentType, filename)
	if err != nil {
		return nil, err
	}

	imageURL, err := s.imageStorage.Save(userID, resolvedType, r, size)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			return nil, err
		}
		return nil, &ValidationError{Message: err.Error()}
	}

	row, err := s.queries.CreateTweetWithImage(ctx, db.CreateTweetWithImageParams{
		UserID:  userID,
		Content: content,
		ImageUrl: sql.NullString{
			String: imageURL,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, &ServiceError{Message: "ツイートの投稿に失敗しました"}
	}

	tweet := tweetFromCreateWithImageRow(row)
	return &tweet, nil
}

func resolveImageContentType(contentType, filename string) (string, error) {
	if contentType != "" {
		if _, ok := storage.AllowedImageMIME(contentType); ok {
			return contentType, nil
		}
	}

	switch strings.ToLower(filepath.Ext(filename)) {
	case ".jpg", ".jpeg":
		return "image/jpeg", nil
	case ".png":
		return "image/png", nil
	case ".gif":
		return "image/gif", nil
	case ".webp":
		return "image/webp", nil
	}

	return "", &ValidationError{Message: "対応していない画像形式です（JPEG, PNG, GIF, WebP のみ）"}
}

func (s *tweetService) GetUserTweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]TweetDetail, error) {
	cursorValue := int32(0)
	if cursor != nil {
		cursorValue = *cursor
	}

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
		base := tweetFromCursorRow(t)
		details = append(details, TweetDetail{Tweet: base, LikeCount: n, RetweetCount: rt})
	}
	return details, nil
}

func (s *tweetService) GetTweetByID(ctx context.Context, id int32) (*TweetDetail, error) {
	tweet, err := s.queries.GetTweetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Tweetが見つかりませんでした")
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

	base := tweetFromGetByIDRow(tweet)
	return &TweetDetail{Tweet: base, LikeCount: likeCount, RetweetCount: retweetCount}, nil
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

	likeByTweetID := make(map[int32]int64, len(likeRows))
	for _, row := range likeRows {
		likeByTweetID[row.TweetID] = row.LikeCount
	}
	retweetByTweetID := make(map[int32]int64, len(retweetRows))
	for _, row := range retweetRows {
		retweetByTweetID[row.TweetID] = row.RetweetCount
	}

	out := make(map[int32]TweetDetail, len(tweets))
	for _, tweet := range tweets {
		base := tweetFromIDsRow(tweet)
		out[tweet.ID] = TweetDetail{
			Tweet:        base,
			LikeCount:    likeByTweetID[tweet.ID],
			RetweetCount: retweetByTweetID[tweet.ID],
		}
	}
	return out, nil
}
