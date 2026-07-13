package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthRequired はログインが必要なエンドポイントを保護するミドルウェア
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// Redisに保存されているユーザーIDを取得
		userID := session.Get("user_id")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
			c.Abort()
			return
		}

		// ユーザーIDをコンテキストに保存
		c.Set("user_id", userID)
		c.Next()
	}
}

// GetUserID はコンテキストからユーザーIDを取得するヘルパー関数
func GetUserID(c *gin.Context) (int32, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("ユーザーIDが見つかりません")
	}

	// ユーザーIDをint32に変換
	userID, ok := userIDInterface.(int32)
	if !ok {
		return 0, errors.New("ユーザーIDの型変換に失敗しました")
	}

	return userID, nil
}
