package services

import (
	"context"
	"database/sql"
	"errors"
	db "golang_twitter/db/sqlc"
	"golang_twitter/infrastructure/storage"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

// ErrUserProfileNotFound はプロフィールが未作成の場合に返す
var ErrUserProfileNotFound = errors.New("プロフィールが見つかりません")

type UserProfileService interface {
	CreateUserProfile(ctx context.Context, userID int32, name, bio, imageURL, location string) (*db.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID int32) (*db.UserProfile, error)
	UpdateUserProfile(ctx context.Context, userID int32, name, bio, imageURL, location string) (*db.UserProfile, error)
	PresignImageUpload(ctx context.Context, userID int32, contentType string, size int64) (*storage.PresignUploadResult, error)
	CompleteImageUpload(ctx context.Context, userID int32, key string) (*db.UserProfile, error)
}

type userProfileService struct {
	queries      *db.Queries
	imageStorage storage.ImageStorage
}

func NewUserProfileService(queries *db.Queries, imageStorage storage.ImageStorage) UserProfileService {
	return &userProfileService{
		queries:      queries,
		imageStorage: imageStorage,
	}
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

// PresignImageUpload はプロフィール画像を S3 へ直接アップロードするための URL と key を発行する
func (s *userProfileService) PresignImageUpload(ctx context.Context, userID int32, contentType string, size int64) (*storage.PresignUploadResult, error) {
	if s.imageStorage == nil {
		return nil, &ServiceError{Message: "画像ストレージが設定されていません"}
	}

	result, err := s.imageStorage.PresignUpload(userID, contentType, size)
	if err != nil {
		return nil, &ValidationError{Message: err.Error()}
	}

	return result, nil
}

// CompleteImageUpload はアップロード済みの画像を確認し、プロフィール画像として保存する
func (s *userProfileService) CompleteImageUpload(ctx context.Context, userID int32, key string) (*db.UserProfile, error) {
	if s.imageStorage == nil {
		return nil, &ServiceError{Message: "画像ストレージが設定されていません"}
	}

	if err := storage.ValidateKeyForUser(userID, key); err != nil {
		return nil, &ValidationError{Message: err.Error()}
	}

	imageURL, err := s.imageStorage.ConfirmUpload(key)
	if err != nil {
		return nil, &ValidationError{Message: err.Error()}
	}

	profile, err := s.queries.UpdateUserProfileImage(ctx, db.UpdateUserProfileImageParams{
		UserID:   userID,
		ImageUrl: imageURL,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserProfileNotFound
		}
		return nil, &ServiceError{Message: "プロフィール画像の保存に失敗しました"}
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
