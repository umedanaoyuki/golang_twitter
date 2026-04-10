package services

import (
	"context"
	db "golang_twitter/db/sqlc"
	"golang_twitter/mailer"
	"golang_twitter/messages"
	"log"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email, password string) (*db.User, error)
}

type authService struct {
	queries *db.Queries
	mailer  mailer.Mailer
}

func NewAuthService(queries *db.Queries, mailer mailer.Mailer) AuthService {
	return &authService{
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

	// 4. メール送信
	if err := s.mailer.SendWelcomeEmail(user.Email); err != nil {
		log.Printf("[ERROR] ウェルカムメール送信失敗 - ユーザー: %s, エラー: %v", user.Email, err)
	}

	return &user, nil
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
