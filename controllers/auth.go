package controllers

import (
	"golang_twitter/services" // Serviceパッケージをインポート
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
    authService services.AuthService
}

func NewAuthController(s services.AuthService) *AuthController {
    return &AuthController{authService: s}
}

type RegisterInput struct {
    MailAddress string `json:"mail_address" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=8,max=15"`
}

func (ctrl *AuthController) Register(c *gin.Context) {
    var input RegisterInput

    // 1. 入力の形式チェック (JSONの構造自体)
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2. Serviceを呼び出してビジネスロジックを実行
    user, err := ctrl.authService.RegisterUser(c.Request.Context(), input.MailAddress, input.Password)
    if err != nil {
        // エラー内容に応じてステータスコードを分岐（簡易版）
        if err.Error() == "パスワードは..." { // 実際は独自のError型で判定するのが理想
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusConflict, gin.H{"error": "登録に失敗しました"})
        }
        return
    }

    // 3. レスポンスの返却
    c.JSON(http.StatusCreated, gin.H{
        "message": "ユーザー登録が完了しました",
        "user": gin.H{
            "id":         user.ID,
            "email":      user.Email,
            "created_at": user.CreatedAt,
        },
    })
} 