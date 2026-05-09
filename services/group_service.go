package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

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
		if err == sql.ErrNoRows {
			return nil
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