package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"

	"github.com/lib/pq"
)

const (
	pgUniqueViolation    = "23505"
	pgForeignKeyViolation = "23503"
)

type GroupService interface {
	CreateGroup(ctx context.Context, userID int32, name string, memberUserIDs []int32) (*db.Group, error)
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

func (s *groupService) CreateGroup(ctx context.Context, creatorID int32, name string, memberUserIDs []int32) (*db.Group, error) {
	memberIDs := uniqueMemberIDs(creatorID, memberUserIDs)

	// Twiterにすべてのメンバーが登録されているか？（正しいuserIDかのチェック）
	for _, memberID := range memberIDs {
		if _, err := s.queries.GetUserDetailByUserID(ctx, memberID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, &ValidationError{Message: "存在しないユーザーが含まれています"}
			}
			return nil, err
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, &ServiceError{Message: "トランザクションの開始に失敗しました"}
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	group, err := qtx.CreateGroup(ctx, db.CreateGroupParams{
		Name:   name,
		UserID: creatorID,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgUniqueViolation {
			return nil, &ValidationError{Message: "同じ名前のグループが既に存在します"}
		}
		return nil, err
	}

	for _, memberID := range memberIDs {
		if _, err := qtx.CreateGroupMember(ctx, db.CreateGroupMemberParams{
			GroupID: group.ID,
			UserID:  memberID,
		}); err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == pgForeignKeyViolation {
				return nil, &ValidationError{Message: "存在しないユーザーが含まれています"}
			}
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, &ServiceError{Message: "グループの作成に失敗しました"}
	}

	return &group, nil
}

func (s *groupService) GetGroups(ctx context.Context, userID int32) ([]db.Group, error) {
	groups, err := s.queries.GetGroupsByMemberUserID(ctx, userID)
	if err != nil {
		return nil, &ServiceError{Message: "グループの取得に失敗しました"}
	}
	return groups, nil
}

func uniqueMemberIDs(creatorID int32, memberUserIDs []int32) []int32 {
	seen := map[int32]struct{}{creatorID: {}}
	ids := []int32{creatorID}

	for _, id := range memberUserIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	return ids
}
