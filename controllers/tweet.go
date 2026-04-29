package controllers

import (
	"net/http"

	"golang_twitter/middleware"
	"golang_twitter/services"

	"github.com/gin-gonic/gin"
)

type TweetController struct {
	tweetService services.TweetService
}

func NewTweetController(tweetService services.TweetService) *TweetController {
	return &TweetController{
		tweetService: tweetService,
	}
}

type CreateTweetInput struct {
	Content string `json:"content" binding:"required"`
}

type GetUserTweetsUri struct {
	UserID int32 `uri:"user_id" binding:"required,min=1"`
}

type GetUserTweetsQuery struct {
	Cursor *int32 `form:"cursor" binding:"omitempty,min=1"`
	Limit  int32  `form:"limit" binding:"omitempty,min=1,max=100"`
}

// 投稿
func (ctrl *TweetController) CreateTweet(c *gin.Context) {
	// ミドルウェアからユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var input CreateTweetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tweetの作成
	tweet, err := ctrl.tweetService.CreateTweet(c.Request.Context(), userID, input.Content)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tweet": gin.H{
			"user_id":    tweet.UserID,
			"content":    tweet.Content,
		},
	})
}

// user_idのツイート一覧を取得
func (ctrl *TweetController) GetUserTweets(c *gin.Context) {
	// パスパラメータのバリデーション
	var uriParams GetUserTweetsUri
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// クエリパラメータのバリデーション（デフォルト値を設定）
	var queryParams GetUserTweetsQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// limitのデフォルト値を設定
	limit := queryParams.Limit
	if limit == 0 {
		limit = 20
	}

	// ツイート一覧を取得
	tweets, err := ctrl.tweetService.GetUserTweetsWithCursor(
		c.Request.Context(), 
		uriParams.UserID, 
		queryParams.Cursor, 
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 次のカーソルを計算
	var nextCursor *int32
	if len(tweets) > 0 {
		lastTweetID := tweets[len(tweets)-1].ID
		nextCursor = &lastTweetID
	}

	c.JSON(http.StatusOK, gin.H{
		"tweets":      tweets,
		"next_cursor": nextCursor,
		"has_more":    len(tweets) == int(limit),
	})
}