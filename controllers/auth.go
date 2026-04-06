package controllers

import (
	"golang_twitter/messages"
	"golang_twitter/services"
	"net/http"

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