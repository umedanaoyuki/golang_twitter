package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"
)

type UserService interface {
	GetUserDetailByUserID(ctx context.Context, userID int32) (*db.GetUserDetailByUserIDRow, error)
}

type userService struct {
	queries *db.Queries
}

func NewUserService(queries *db.Queries) UserService {
	return &userService{queries: queries}
}

func (s *userService) GetUserDetailByUserID(ctx context.Context, userID int32) (*db.GetUserDetailByUserIDRow, error) {
	user, err := s.queries.GetUserDetailByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ユーザーが見つかりません")
		}
		return nil, &ServiceError{Message: "ユーザー情報の取得に失敗しました"}
	}
	return &user, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID int32) error {
	err := s.queries.DeleteUser(ctx, userID)
	if err != nil {
		return &ServiceError{Message: "ユーザーの削除に失敗しました"}
	}
	return nil
}