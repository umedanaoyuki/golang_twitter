package services

import (
	"context"
	"database/sql"
	db "golang_twitter/db/sqlc"
)

type GroupService interface {
	CreateGroup(ctx context.Context, userID int32) error
	GetGroups(ctx context.Context, userID int32) error
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

func (s *groupService) CreateGroup(ctx context.Context, userID int32) error {
	return nil;
}

func (s *groupService) GetGroups(ctx context.Context, userID int32) error {
	return nil;
}