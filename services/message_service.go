package services

import (
	"context"
	"database/sql"

	db "golang_twitter/db/sqlc"
)

type MessageService interface {
	CreateMessage(ctx context.Context, userID int32, groupID int32, content string) (*db.Message, error)
	GetMessages(ctx context.Context, groupID int32) error
}

type messageService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewMessageService(db *sql.DB, queries *db.Queries) MessageService {
	return &messageService{
		db:      db,
		queries: queries,
	}
}

func (s *messageService) CreateMessage(ctx context.Context, userID int32, groupID int32, content string) (*db.Message, error) {
	message, err := s.queries.CreateMessage(ctx, db.CreateMessageParams{
		UserID:  userID,
		GroupID: groupID,
		Content: content,
	})
	if err != nil {
		return nil, &ServiceError{Message: "メッセージの作成に失敗しました"}
	}
	return &message, nil
}

func (s *messageService) GetMessages(ctx context.Context, groupID int32) error {
	return nil;
}