package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type RetweetService interface {
	CreateRetweet(ctx context.Context, userID int32, tweetID int32) error
	DeleteRetweet(ctx context.Context, userID int32, tweetID int32) error
	GetUserRetweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]UserRetweetItem, error)
}

// UserRetweetItem は API 用の1件分（ツイート詳細のみ JSON に載せる）
type UserRetweetItem struct {
	// ページネーション用。レスポンスには含めない
	RetweetRowID int32 `json:"-"`
	Tweet        *TweetDetail `json:"tweet,omitempty"`
}

type retweetService struct {
	db           *sql.DB
	queries      *db.Queries
	tweetService TweetService
}

func NewRetweetService(db *sql.DB, queries *db.Queries, tweetService TweetService) RetweetService {
	return &retweetService{
		db:           db,
		queries:      queries,
		tweetService: tweetService,
	}
}

func (s *retweetService) CreateRetweet(ctx context.Context, userID int32, tweetID int32) error {
	_, err := s.queries.CreateRetweet(ctx, db.CreateRetweetParams{
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

func (s *retweetService) DeleteRetweet(ctx context.Context, userID int32, tweetID int32) error {
	return s.queries.DeleteRetweet(ctx, db.DeleteRetweetParams{
		UserID:  userID,
		TweetID: tweetID,
	})
}

// ユーザーのリツイート一覧を取得（ページネーションつき）。ツイート本文などは TweetService 経由で一括取得して紐づける
func (s *retweetService) GetUserRetweetsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]UserRetweetItem, error) {
	cursorValue := int32(0)
	if cursor != nil {
		cursorValue = *cursor
	}
	retweets, err := s.queries.GetUserRetweetsWithCursor(ctx, db.GetUserRetweetsWithCursorParams{
		UserID:  userID,
		Column2: cursorValue,
		Limit:   limit,
	})
	if err != nil {
		return nil, err
	}
	if len(retweets) == 0 {
		return []UserRetweetItem{}, nil
	}

	tweetIDs := make([]int32, len(retweets))
	for i := range retweets {
		tweetIDs[i] = retweets[i].TweetID
	}
	detailsByID, err := s.tweetService.GetTweetDetailsByIDs(ctx, tweetIDs)
	if err != nil {
		return nil, err
	}

	items := make([]UserRetweetItem, 0, len(retweets))
	for _, rt := range retweets {
		item := UserRetweetItem{RetweetRowID: rt.ID}
		if d, ok := detailsByID[rt.TweetID]; ok {
			dCopy := d
			item.Tweet = &dCopy
		}
		items = append(items, item)
	}
	return items, nil
}