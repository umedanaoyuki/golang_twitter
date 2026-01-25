package controllers

import (
	db "golang_twitter/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	MailAddress string `json:"mail_address" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context, queries *db.Queries) {
	var input RegisterInput

	// JSONバリデーション
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
		return
	}

	// データベースにユーザーを作成
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