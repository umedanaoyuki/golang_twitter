package main

import (
	"context"
	"database/sql"
	"golang_twitter/controllers"
	db "golang_twitter/db/sqlc"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// 環境変数からDATABASE_URLを取得
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL が設定されていません")
	}

	// データベース接続
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("データベース接続エラー:", err)
	}
	defer conn.Close()

	// 接続確認
	if err := conn.Ping(); err != nil {
		log.Fatal("データベース接続確認エラー:", err)
	}

	log.Println("データベース接続成功")

	ctx := context.Background()
  // 生SQLを事前にDBに送信してチェック
	queries, err := db.Prepare(ctx, conn)
	if err != nil {
		log.Fatal("プリペアードステートメント作成エラー:", err)
	}
	defer queries.Close()

	// Ginルーター設定
	router := gin.Default()

	router.GET("/health_check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.POST("/register", func(c *gin.Context) {
    // クエリを渡す
		controllers.Register(c, queries)
	})

	log.Println("🚀 サーバー起動: http://localhost:8080")
	router.Run()
}