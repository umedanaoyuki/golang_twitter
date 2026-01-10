package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
  // loggerとrecoveryミドルウェア付きGinルーター作成
  r := gin.Default()

  // 簡単なGETエンドポイント定義
  r.GET("/ping", func(c *gin.Context) {
    // JSONレスポンスを返す
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  // ポート8080でサーバー起動（デフォルト）
  // 0.0.0.0:8080（Windowsではlocalhost:8080）で待機
  r.Run()
}