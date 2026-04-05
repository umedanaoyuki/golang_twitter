package services

import (
	"context"
	db "golang_twitter/db/sqlc"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email, password string) (*db.User, error)
}

type authService struct {
	queries *db.Queries
}

func NewAuthService(queries *db.Queries) AuthService {
	return &authService{queries: queries}
}

// RegisterUser はユーザー登録のビジネスロジックを実行
func (s *authService) RegisterUser(ctx context.Context, email, password string) (*db.User, error) {
	// 1. パスワードのバリデーション
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	// 2. パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &ServiceError{Message: "パスワードのハッシュ化に失敗しました"}
	}

	// 3. ユーザーを作成
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, &ServiceError{Message: "このメールアドレスは既に登録されています"}
	}

	return &user, nil
}

func validatePassword(password string) error {
	// 1. 長さチェック（8文字以上15文字以下）
	if len(password) > 15 || len(password) < 8 {
		return &ValidationError{Message: "パスワードは8文字以上15文字以下である必要があります"}
	}

	// 2. 半角英数字のみ + 許可された記号かチェック
	validChars := regexp.MustCompile(`^[a-zA-Z0-9!?_-]+$`)
	if !validChars.MatchString(password) {
		return &ValidationError{Message: "パスワードは半角英数字と記号(!?-_)のみ使用できます"}
	}

	// 3. 小文字が含まれているかチェック
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return &ValidationError{Message: "パスワードには小文字を含める必要があります"}
	}

	// 4. 大文字が含まれているかチェック
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return &ValidationError{Message: "パスワードには大文字を含める必要があります"}
	}

	// 5. 数字が含まれているかチェック
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return &ValidationError{Message: "パスワードには数字を含める必要があります"}
	}

	// 6. 指定された記号が含まれているかチェック
	hasSymbol := strings.ContainsAny(password, "!?-_")
	if !hasSymbol {
		return &ValidationError{Message: "パスワードには記号(!?-_)を1文字以上含める必要があります"}
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
