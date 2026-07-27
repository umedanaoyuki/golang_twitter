package controllers

import (
	"golang_twitter/messages"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8,max=15" example:"Password1!"`
}

// Register godoc
// @Summary      ユーザー登録
// @Description  メールアドレスとパスワードで新規ユーザーを登録する
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      RegisterInput  true  "登録情報"
// @Success      201    {object}  RegisterResponse
// @Failure      400    {object}  ErrorResponse
// @Router       /register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var input RegisterInput

	// 1. JSONの形式チェック
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. ビジネスロジック実行
	user, err := ctrl.authService.RegisterUser(c.Request.Context(), input.Email, input.Password)
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

// ActivateUser godoc
// @Summary      アカウント有効化
// @Description  メールで送付されたトークンでアカウントを有効化する
// @Tags         auth
// @Produce      json
// @Param        token  query     string  true  "アクティベーショントークン"
// @Success      200    {object}  MessageResponse
// @Failure      400    {object}  ErrorResponse
// @Router       /activate [get]
func (ctrl *AuthController) ActivateUser(c *gin.Context) {
	// 1. クエリパラメータからトークンを取得
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "トークンが指定されていません"})
		return
	}

	// 2. アクティベーション処理
	if err := ctrl.authService.ActivateUser(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgUserActivated,
	})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"Password1!"`
}

// Login godoc
// @Summary      ログイン
// @Description  メールアドレスとパスワードでログインし、セッション Cookie を発行する
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      LoginInput  true  "ログイン情報"
// @Success      200    {object}  LoginResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var input LoginInput

	// 1. JSONの形式チェック
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. ログイン処理
	user, err := ctrl.authService.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 3. セッションにユーザーIDを保存
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrSessionSaveFailed})
		return
	}

	// 4. 成功レスポンス（パスワードは返さない）
	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgLoginSuccess,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
		},
	})
}

// Logout godoc
// @Summary      ログアウト
// @Description  セッションを破棄し、セッション Cookie を削除する
// @Tags         auth
// @Produce      json
// @Security     SessionAuth
// @Success      200  {object}  MessageResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
	session := sessions.Default(c)

	// 1. セッションの中身を削除
	session.Clear()

	// 2. MaxAgeを-1にすることでブラウザ側のCookieも削除する
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1,
	})

	// 3. Redisからセッションを削除し、Set-Cookieを発行
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrSessionSaveFailed})
		return
	}

	// 4. 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgLogoutSuccess,
	})
}