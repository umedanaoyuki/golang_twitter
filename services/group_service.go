package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"

	"github.com/lib/pq"
)

const pgUniqueViolation = "23505"

type GroupService interface {
	CreateGroup(ctx context.Context, userID int32, name string) error
	GetGroups(ctx context.Context, userID int32) ([]db.Group, error)
}

type groupService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewGroupService(db *sql.DB, queries *db.Queries) GroupService {
	return &groupService{
		db:      db,
		queries: queries,
	}
}

func (s *groupService) CreateGroup(ctx context.Context, userID int32, name string) error {
	_, err := s.queries.CreateGroup(ctx, db.CreateGroupParams{
		Name:   name,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgUniqueViolation {
			return &ValidationError{Message: "同じ名前のグループが既に存在します"}
		}
		return err
	}
	return nil
}

func (s *groupService) GetGroups(ctx context.Context, userID int32) ([]db.Group, error) {
	groups, err := s.queries.GetGroupsByUserID(ctx, userID)
	if err != nil {
		return nil, &ServiceError{Message: "グループの取得に失敗しました"}
	}
	return groups, nil
}