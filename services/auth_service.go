package services

import (
	"context"
	"errors"
	db "golang_twitter/db/sqlc"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// AuthService のインターフェース定義
type AuthService interface {
    RegisterUser(ctx context.Context, email, password string) (db.User, error)
}

type authService struct {
    queries *db.Queries
}

func NewAuthService(queries *db.Queries) AuthService {
    return &authService{queries: queries}
}

// RegisterUser ビジネスロジックの本体
func (s *authService) RegisterUser(ctx context.Context, email, password string) (db.User, error) {
    // 1. パスワードのバリデーション
    if err := validatePassword(password); err != nil {
        return db.User{}, err
    }

    // 2. パスワードハッシュ化
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return db.User{}, errors.New("パスワードのハッシュ化に失敗しました")
    }

    // 3. ユーザー作成（Repository/sqlcの呼び出し）
    user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
        Email:    email,
        Password: string(hashedPassword),
    })
    if err != nil {
        return db.User{}, err // 重複エラーなどはControllerで判定
    }

    return user, nil
}

// 内部ロジックとしてのバリデーション（外部に公開しない場合は小文字開始）
func validatePassword(password string) error {
    if len(password) < 8 || len(password) > 15 {
        return errors.New("パスワードは8文字以上15文字以下である必要があります")
    }
    validChars := regexp.MustCompile(`^[a-zA-Z0-9!?_-]+$`)
    if !validChars.MatchString(password) {
        return errors.New("パスワードは半角英数字と記号(!?-_)のみ使用できます")
    }
    // ...（他のバリデーションもここに移動）...
    hasSymbol := strings.ContainsAny(password, "!?-_")
    if !hasSymbol {
        return errors.New("パスワードには記号(!?-_)を1文字以上含める必要があります")
    }
    return nil
}