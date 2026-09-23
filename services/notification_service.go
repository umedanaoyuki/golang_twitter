package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"
	"log"
	"time"
)

const (
	NotificationTypeLike    = "like"
	NotificationTypeFollow  = "follow"
	NotificationTypeComment = "comment"
)

type NotificationService interface {
	// NotifyLike はツイートにいいねされたことをツイート作者に通知する
	NotifyLike(ctx context.Context, actorID int32, tweetID int32)
	// NotifyFollow はフォローされたことをフォローされたユーザーに通知する
	NotifyFollow(ctx context.Context, actorID int32, followedUserID int32)
	// NotifyComment はツイートにコメントされたことをツイート作者に通知する
	NotifyComment(ctx context.Context, actorID int32, tweetID int32, commentID int32)
	// GetNotificationsWithCursor はログインユーザー宛の通知一覧をカーソルページネーションで取得する
	GetNotificationsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]NotificationItem, error)
}

// NotificationItem は API 用の通知1件分
type NotificationItem struct {
	ID        int32      `json:"id"`
	UserID    int32      `json:"user_id"`
	ActorID   int32      `json:"actor_id"`
	Type      string     `json:"type"`
	TweetID   *int32     `json:"tweet_id"`
	CommentID *int32     `json:"comment_id"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type notificationService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewNotificationService(db *sql.DB, queries *db.Queries) NotificationService {
	return &notificationService{
		db:      db,
		queries: queries,
	}
}

// 通知の作成は本体処理（いいね等）のおまけ扱いとし、失敗してもログに残すだけで本体処理は失敗させない
func (s *notificationService) create(ctx context.Context, params db.CreateNotificationParams) {
	if _, err := s.queries.CreateNotification(ctx, params); err != nil {
		log.Printf("通知の作成に失敗しました: type=%s user_id=%d actor_id=%d err=%v", params.Type, params.UserID, params.ActorID, err)
	}
}

// ツイート作者を取得する。ツイートが存在しない場合は ok=false
func (s *notificationService) tweetOwner(ctx context.Context, tweetID int32) (int32, bool) {
	tweet, err := s.queries.GetTweetByID(ctx, tweetID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("通知対象ツイートの取得に失敗しました: tweet_id=%d err=%v", tweetID, err)
		}
		return 0, false
	}
	return tweet.UserID, true
}

func (s *notificationService) NotifyLike(ctx context.Context, actorID int32, tweetID int32) {
	ownerID, ok := s.tweetOwner(ctx, tweetID)
	// 自分のツイートへのいいねは通知しない
	if !ok || ownerID == actorID {
		return
	}
	s.create(ctx, db.CreateNotificationParams{
		UserID:  ownerID,
		ActorID: actorID,
		Type:    NotificationTypeLike,
		TweetID: sql.NullInt32{Int32: tweetID, Valid: true},
	})
}

func (s *notificationService) NotifyFollow(ctx context.Context, actorID int32, followedUserID int32) {
	if actorID == followedUserID {
		return
	}
	s.create(ctx, db.CreateNotificationParams{
		UserID:  followedUserID,
		ActorID: actorID,
		Type:    NotificationTypeFollow,
	})
}

func (s *notificationService) NotifyComment(ctx context.Context, actorID int32, tweetID int32, commentID int32) {
	ownerID, ok := s.tweetOwner(ctx, tweetID)
	// 自分のツイートへのコメントは通知しない
	if !ok || ownerID == actorID {
		return
	}
	s.create(ctx, db.CreateNotificationParams{
		UserID:    ownerID,
		ActorID:   actorID,
		Type:      NotificationTypeComment,
		TweetID:   sql.NullInt32{Int32: tweetID, Valid: true},
		CommentID: sql.NullInt32{Int32: commentID, Valid: true},
	})
}

func (s *notificationService) GetNotificationsWithCursor(ctx context.Context, userID int32, cursor *int32, limit int32) ([]NotificationItem, error) {
	cursorValue := int32(0)
	if cursor != nil {
		cursorValue = *cursor
	}

	rows, err := s.queries.GetNotificationsByUserIDWithCursor(ctx, db.GetNotificationsByUserIDWithCursorParams{
		UserID:  userID,
		Column2: cursorValue,
		Limit:   limit,
	})
	if err != nil {
		return nil, &ServiceError{Message: "通知の取得に失敗しました"}
	}

	items := make([]NotificationItem, 0, len(rows))
	for _, n := range rows {
		items = append(items, toNotificationItem(n))
	}
	return items, nil
}

// DB モデルの NULL 許容型（sql.NullInt32 等）をポインタに変換して JSON で扱いやすくする
func toNotificationItem(n db.Notification) NotificationItem {
	item := NotificationItem{
		ID:        n.ID,
		UserID:    n.UserID,
		ActorID:   n.ActorID,
		Type:      n.Type,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
	if n.TweetID.Valid {
		v := n.TweetID.Int32
		item.TweetID = &v
	}
	if n.CommentID.Valid {
		v := n.CommentID.Int32
		item.CommentID = &v
	}
	if n.ReadAt.Valid {
		v := n.ReadAt.Time
		item.ReadAt = &v
	}
	return item
}
