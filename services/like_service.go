package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type LikeService interface {
	CreateLike(ctx context.Context, userID int32, tweetID int32) error
	DeleteLike(ctx context.Context, userID int32, tweetID int32) error
	GetUserLikesWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]UserLikeItem, error)
}

// UserLikeItem は API 用の1件分（ツイート詳細のみ JSON に載せる）
type UserLikeItem struct {
	// ページネーション用。レスポンスには含めない
	LikeRowID int32        `json:"-"`
	Tweet     *TweetDetail `json:"tweet,omitempty"`
}

type likeService struct {
	db           *sql.DB
	queries      *db.Queries
	tweetService TweetService
}

func NewLikeService(db *sql.DB, queries *db.Queries, tweetService TweetService) LikeService {
	return &likeService{
		db:           db,
		queries:      queries,
		tweetService: tweetService,
	}
}

func (s *likeService) CreateLike(ctx context.Context, userID int32, tweetID int32) error {
	_, err := s.queries.CreateLike(ctx, db.CreateLikeParams{
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

func (s *likeService) DeleteLike(ctx context.Context, userID int32, tweetID int32) error {
	return s.queries.DeleteLike(ctx, db.DeleteLikeParams{
		UserID:  userID,
		TweetID: tweetID,
	})
}

// ユーザーのいいね一覧を取得（ページネーションつき）。ツイート本文などは TweetService 経由で一括取得して紐づける
func (s *likeService) GetUserLikesWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]UserLikeItem, error) {
	cursorValue := int32(0)
	if cursor != nil {
		cursorValue = *cursor
	}
	likes, err := s.queries.GetUserLikesWithCursor(ctx, db.GetUserLikesWithCursorParams{
		UserID:  userID,
		Column2: cursorValue,
		Limit:   limit,
	})
	if err != nil {
		return nil, err
	}
	if len(likes) == 0 {
		return []UserLikeItem{}, nil
	}

	tweetIDs := make([]int32, len(likes))
	for i := range likes {
		tweetIDs[i] = likes[i].TweetID
	}
	detailsByID, err := s.tweetService.GetTweetDetailsByIDs(ctx, tweetIDs)
	if err != nil {
		return nil, err
	}

	items := make([]UserLikeItem, 0, len(likes))
	for _, like := range likes {
		item := UserLikeItem{LikeRowID: like.ID}
		if d, ok := detailsByID[like.TweetID]; ok {
			dCopy := d
			item.Tweet = &dCopy
		}
		items = append(items, item)
	}
	return items, nil
}