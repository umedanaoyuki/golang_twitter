package main

import (
	"context"
	"database/sql"
	"golang_twitter/controllers"
	db "golang_twitter/db/sqlc"
	"golang_twitter/infrastructure/email/mailcatcher"
	"golang_twitter/mailer"
	"golang_twitter/services"
	"golang_twitter/sessionInfo"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var LoginInfo sessionInfo.SessionInfo

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

	var emailMailer mailer.Mailer
	if os.Getenv("ENV") == "development" {
		emailMailer = mailcatcher.NewMailCatcherMailer()
	} else {
		// TODO: SMTPの設定値を書く（環境変数）
		// emailMailer = smtp.NewEmailSender()
	}

	// サービスの初期化
	authService := services.NewAuthService(conn, queries, emailMailer)
	tweetService := services.NewTweetService(conn, queries)

	// コントローラーの初期化
	authController := controllers.NewAuthController(authService)
	tweetController := controllers.NewTweetController(tweetService)

	// Ginルーター設定
	router := gin.Default()

	router.GET("/health_check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// ユーザー登録
	// セッションの設定
	store := cookie.NewStore([]byte("secret"))
	id, err := uuid.NewV7()
	// keyは一意な値(uuid)
	router.Use(sessions.Sessions(id.String(), store))

	// ユーザー登録エンドポイント
	router.POST("/register", authController.Register)

	// ユーザーアクティベーション
	router.GET("/activate", authController.ActivateUser)

	// ログイン
	router.POST("/login", authController.Login)

	// ツイート関連
	router.POST("/tweets", tweetController.CreateTweet)
	
	log.Println("サーバー起動: http://localhost:8080")
	router.Run()
}

// NOTE: セッションをチェックするために後ほど使用？
func sessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)
		LoginInfo.UserId = session.Get("UserId")

		// セッションがない場合、ログインフォームをだす
		if LoginInfo.UserId == nil {
			log.Println("ログインしていません")
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort() // これがないと続けて処理されてしまう
		} else {
			c.Set("UserId", LoginInfo.UserId) // ユーザidをセット
			c.Next()
		}
		log.Println("ログインチェック終わり")
	}
}