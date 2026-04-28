package services

import (
	"context"
	db "golang_twitter/db/sqlc"
)

type UserService interface {
	GetUserDetailByUserID(ctx context.Context, userID int32) (*db.User, error)
}

type userService struct {
	queries *db.Queries
}

func NewUserService(queries *db.Queries) UserService {
	return &userService{queries: queries}
}

func (s *userService) GetUserDetailByUserID(ctx context.Context, userID int32) (*db.User, error) {
	user, err := s.queries.GetUserDetailByUserID(ctx, userID)
	if err != nil {
		return nil, &ServiceError{Message: "ユーザー情報の取得に失敗しました"}
	}
	return &user, nil
}