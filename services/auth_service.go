package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	db "golang_twitter/db/sqlc"
	"golang_twitter/mailer"
	"golang_twitter/messages"
	"log"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email, password string) (*db.User, error)
	ActivateUser(ctx context.Context, token string) error
	Login(ctx context.Context, email, password string) (*db.User, error)
}

type authService struct {
	db      *sql.DB
	queries *db.Queries
	mailer  mailer.Mailer
}

func NewAuthService(db *sql.DB, queries *db.Queries, mailer mailer.Mailer) AuthService {
	return &authService{
		db:      db,
		queries: queries,
		mailer:  mailer,
	}
}

func (s *authService) RegisterUser(ctx context.Context, email, password string) (*db.User, error) {
	// 1. パスワードのバリデーション
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	// 2. パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &ServiceError{Message: messages.ErrPasswordHashFailed}
	}

	// 3. ユーザーを作成
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:    email,
		Password: string(hashedPassword),
		IsActive: false,
	})
	if err != nil {
		return nil, &ServiceError{Message: messages.ErrEmailAlreadyExists}
	}

	// 4. アクティベーショントークンを生成
	token := generateRandomToken()
	expiredAt := time.Now().Add(24 * time.Hour)

	// 5. トークンをDBに保存
	_, err = s.queries.CreateUserActivation(ctx, db.CreateUserActivationParams{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: expiredAt,
	})
	if err != nil {
		log.Printf("[ERROR] アクティベーショントークン作成失敗 - ユーザー: %s, エラー: %v", user.Email, err)
	}

	// 6. ウェルカムメール送信（トークン付き）
	if err := s.mailer.SendWelcomeEmail(user.Email, token); err != nil {
		log.Printf("[ERROR] ウェルカムメール送信失敗 - ユーザー: %s, エラー: %v", user.Email, err)
	}

	return &user, nil
}

// ActivateUser はトークンを使用してユーザーをアクティブ化
func (s *authService) ActivateUser(ctx context.Context, token string) error {
	// 1. トークンを検証（有効期限チェックも含む）
	activation, err := s.queries.GetUserActivationByToken(ctx, token)
	if err != nil {
		return &ServiceError{Message: "無効なトークンまたは期限切れです"}
	}

	// 2. トランザクション開始
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return &ServiceError{Message: "トランザクション開始に失敗しました"}
	}
	// ロールバックの予約
	defer tx.Rollback()

	// 3. トランザクション用のQueriesを作成
	qtx := s.queries.WithTx(tx)

	// 4. ユーザーをアクティブ化
	if err := qtx.UpdateUserIsActive(ctx, db.UpdateUserIsActiveParams{
		ID:       activation.UserID,
		IsActive: true,
	}); err != nil {
		return &ServiceError{Message: "ユーザーのアクティブ化に失敗しました"}
	}

	// 5. トークンを削除（使用済み）
	if err := qtx.DeleteUserActivation(ctx, token); err != nil {
		return &ServiceError{Message: "トークン削除に失敗しました"}
	}

	return tx.Commit()
}

// Login はメールアドレスとパスワードでログイン
func (s *authService) Login(ctx context.Context, email, password string) (*db.User, error) {
	// 1. メールアドレスでユーザーを検索
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, &ServiceError{Message: "メールアドレスまたはパスワードが正しくありません"}
	}

	// 2. アカウントがアクティブ化されているかチェック
	if !user.IsActive {
		return nil, &ServiceError{Message: "アカウントが有効化されていません"}
	}

	// 3. パスワードを照合
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, &ServiceError{Message: "メールアドレスまたはパスワードが正しくありません"}
	}

	// 4. 認証成功（パスワードは返さない）
	return &user, nil
}

// generateRandomToken はランダムなトークンを生成
func generateRandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func validatePassword(password string) error {
	// 1. 長さチェック（8文字以上15文字以下）
	if len(password) > 15 || len(password) < 8 {
		return &ValidationError{Message: messages.ErrPasswordLength}
	}

	// 2. 半角英数字のみ + 許可された記号かチェック
	validChars := regexp.MustCompile(`^[a-zA-Z0-9!?_-]+$`)
	if !validChars.MatchString(password) {
		return &ValidationError{Message: messages.ErrPasswordInvalidChars}
	}

	// 3. 小文字が含まれているかチェック
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return &ValidationError{Message: messages.ErrPasswordRequireLower}
	}

	// 4. 大文字が含まれているかチェック
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return &ValidationError{Message: messages.ErrPasswordRequireUpper}
	}

	// 5. 数字が含まれているかチェック
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return &ValidationError{Message: messages.ErrPasswordRequireNumber}
	}

	// 6. 指定された記号が含まれているかチェック
	hasSymbol := strings.ContainsAny(password, "!?-_")
	if !hasSymbol {
		return &ValidationError{Message: messages.ErrPasswordRequireSymbol}
	}

	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
