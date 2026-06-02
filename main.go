package main

import (
	"context"
	"database/sql"
	"golang_twitter/controllers"
	db "golang_twitter/db/sqlc"
	_ "golang_twitter/docs"
	"golang_twitter/infrastructure/email/mailcatcher"
	"golang_twitter/mailer"
	"golang_twitter/middleware"
	"golang_twitter/services"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Golang Twitter API
// @version         1.0
// @description     Twitter クローン API
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey  SessionAuth
// @in                          cookie
// @name                        golang_twitter_session

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
	likeService := services.NewLikeService(conn, queries)
	userService := services.NewUserService(queries)
	bookmarkService := services.NewBookmarkService(conn, queries)
	groupService := services.NewGroupService(conn, queries)
	messageService := services.NewMessageService(conn, queries)
	retweetService := services.NewRetweetService(conn, queries, tweetService)
	followService := services.NewFollowService(conn, queries)

	// コントローラーの初期化
	authController := controllers.NewAuthController(authService)
	tweetController := controllers.NewTweetController(tweetService)
	likeController := controllers.NewLikeController(likeService)
	
	userController := controllers.NewUserController(userService)
	bookmarkController := controllers.NewBookmarkController(bookmarkService)
	groupController := controllers.NewGroupController(groupService)
	messageController := controllers.NewMessageController(messageService)
	retweetController := controllers.NewRetweetController(retweetService)
	followController := controllers.NewFollowController(followService)
	// Ginルーター設定
	router := gin.Default()

	// Redisセッション設定
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis_server:6379"
	}
	
	// Redis Store
	// memo: 第5引数がパスワード（今回は空文字）
	store, err := redis.NewStore(10, "tcp", redisHost, "", "", []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		log.Fatal("Redis接続エラー:", err)
	}
	router.Use(sessions.Sessions("golang_twitter_session", store))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health_check", controllers.HealthCheck)

	// 認証不要のエンドポイント
	// ユーザー登録
	router.POST("/register", authController.Register)
	// ユーザーアクティベーション
	router.GET("/activate", authController.ActivateUser)
	// ログイン
	router.POST("/login", authController.Login)
	// ユーザーのツイート一覧取得
	router.GET("/users/:user_id/tweets", tweetController.GetUserTweets)
	// ユーザーのリツイート一覧取得
	router.GET("/users/:user_id/retweets", retweetController.GetUserRetweets)
	// ユーザー情報取得
	router.GET("/users/:user_id", userController.GetUserByID)

	// 1つのTweet取得
	router.GET("/tweets/:id", tweetController.GetTweetByID)

	// 認証が必要なエンドポイント
	authorized := router.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		// ツイート投稿
		authorized.POST("/tweets", tweetController.CreateTweet)
		authorized.DELETE("/tweets/:id", tweetController.DeleteTweet)

		// リツイート機能
		authorized.POST("/tweets/:id/retweet", retweetController.CreateRetweet)
		authorized.DELETE("/tweets/:id/retweet", retweetController.DeleteRetweet)
		// 退会
		authorized.DELETE("/user", userController.DeleteUser)
		// ブックマーク機能
		authorized.POST("/tweets/:id/bookmark", bookmarkController.CreateBookmark)
		authorized.DELETE("/tweets/:id/bookmark", bookmarkController.DeleteBookmark)
		authorized.GET("/tweets/bookmarks", bookmarkController.GetBookmarksByUserId)
		// いいね機能
		authorized.POST("/tweets/:id/like", likeController.CreateLike)
		authorized.DELETE("/tweets/:id/like", likeController.DeleteLike)

		/* グループ内でのメッセージ機能 */
		// グループ作成
		authorized.POST("/groups", groupController.CreateGroup)
		// グループ一覧取得
		authorized.GET("/groups", groupController.GetGroups)
		// グループ内でのメッセージ送信
		authorized.POST("/groups/:group_id/messages", messageController.CreateMessage)
		// グループ内でのメッセージ一覧取得
		authorized.GET("/groups/:group_id/messages", messageController.GetMessages)

		// ユーザーフォロー機能
		authorized.POST("/users/:user_id/follow", followController.CreateFollow)
		authorized.DELETE("/users/:user_id/follow", followController.DeleteFollow)
		// ユーザーフォロワー一覧取得
		authorized.GET("/users/:user_id/followers", followController.GetFollowersByUserId)
		// ユーザーフォロー中一覧取得
		authorized.GET("/users/:user_id/following", followController.GetFollowingByUserId)
	}

	log.Println("サーバー起動: http://localhost:8080")
	router.Run()
}