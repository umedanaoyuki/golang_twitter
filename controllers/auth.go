package controllers

import (
	db "golang_twitter/db/sqlc"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	MailAddress string `json:"mail_address" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=15"`
}

// パスワードの複雑性をチェックする関数
func validatePassword(password string) error {
	// 1. 長さチェック（8文字以上15文字以下）
	if len(password) < 8 || len(password) > 15 {
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

// カスタムバリデーションエラー型
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func Register(c *gin.Context, queries *db.Queries) {
	var input RegisterInput

	// 基本的なJSONバリデーション
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// パスワードの複雑性チェック
	if err := validatePassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
		return
	}

	// データベースにユーザーを作成（CreateUserはsqlcが自動生成）
	user, err := queries.CreateUser(c.Request.Context(), db.CreateUserParams{
		Email:    input.MailAddress,
		Password: string(hashedPassword),
	})
	if err != nil {
		// メールアドレスが既に存在する場合など
		c.JSON(http.StatusConflict, gin.H{"error": "このメールアドレスは既に登録されています"})
		return
	}

	// 成功レスポンス（パスワードは返さない）
	c.JSON(http.StatusCreated, gin.H{
		"message": "ユーザー登録が完了しました",
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}