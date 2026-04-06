package controllers

import (
	"golang_twitter/infrastructure/email"
	"golang_twitter/messages"
	"golang_twitter/services"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type RegisterInput struct {
	MailAddress string `json:"mail_address" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=15"`
}

type EmailInput struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username"`
}

// ユーザー登録
func (ctrl *AuthController) Register(c *gin.Context) {
	var input RegisterInput

	// 1. JSONの形式チェック
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. ビジネスロジック実行
	user, err := ctrl.authService.RegisterUser(c.Request.Context(), input.MailAddress, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 成功レスポンス（パスワードは返さない）
	c.JSON(http.StatusCreated, gin.H{
		"message": messages.MsgUserRegistered,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}

// ウェルカムメール送信
func (ctrl *AuthController) SendWelcomeEmail(c *gin.Context) {
	var input EmailInput

	// 1. JSONの形式チェック
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Usernameが指定されていない場合はデフォルト値を使用
	username := input.Username
	if username == "" {
		username = "テストユーザー"
	}

	// 2. メール送信処理
	config := email.NewMailCatcherConfig()
	sender := email.NewEmailSender(config)
	templatePath := filepath.Join(email.GetTemplateDir(), "welcome.html")

	data := map[string]interface{}{
		"Username": username,
		"Email":    input.Email,
		"AppName":  "Twitter Clone",
	}

	if err := sender.SendEmail(
		[]string{input.Email},
		"ようこそ Twitter Clone へ",
		templatePath,
		data,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message": "ウェルカムメールを送信しました",
		"email":   input.Email,
	})
}