package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

type UserProfileService interface {
	CreateUserProfile(ctx context.Context, userID int32, name, bio, imageURL, location string) (*db.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID int32) (*db.UserProfile, error)
	UpdateUserProfile(ctx context.Context, userID int32, name, bio, imageURL, location string) (*db.UserProfile, error)
}

type userProfileService struct {
	queries *db.Queries
}

func NewUserProfileService(queries *db.Queries) UserProfileService {
	return &userProfileService{queries: queries}
}

func (s *userProfileService) CreateUserProfile(ctx context.Context, userID int32, name, bio, imageURL, location string) (*db.UserProfile, error) {
	profile, err := s.queries.CreateUserProfile(ctx, db.CreateUserProfileParams{
		UserID:   userID,
		Name:     resolveProfileName(name, userID),
		Bio:      bio,
		ImageUrl: imageURL,
		Location: location,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && string(pqErr.Code) == pgUniqueViolation {
			return nil, &ValidationError{Message: "プロフィールは既に作成されています"}
		}
		return nil, &ServiceError{Message: "プロフィールの作成に失敗しました"}
	}
	return &profile, nil
}

func (s *userProfileService) GetUserProfileByUserID(ctx context.Context, userID int32) (*db.UserProfile, error) {
	profile, err := s.queries.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &ValidationError{Message: "プロフィールが見つかりません"}
		}
		return nil, &ServiceError{Message: "プロフィールの取得に失敗しました"}
	}
	return &profile, nil
}

func (s *userProfileService) UpdateUserProfile(ctx context.Context, userID int32, name, bio, imageURL, location string) (*db.UserProfile, error) {
	profile, err := s.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		UserID:   userID,
		Name:     resolveProfileName(name, userID),
		Bio:      bio,
		ImageUrl: imageURL,
		Location: location,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &ValidationError{Message: "プロフィールが見つかりません"}
		}
		return nil, &ServiceError{Message: "プロフィールの更新に失敗しました"}
	}
	return &profile, nil
}

// resolveProfileName は name が未指定の場合に user_id を名前として使う
func resolveProfileName(name string, userID int32) string {
	if strings.TrimSpace(name) == "" {
		return `ユーザーID_` + strconv.Itoa(int(userID))
	}
	return name
}
