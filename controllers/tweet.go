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

// CreateTweet はツイートを投稿
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

	// ツイートを作成
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
		"message": "ツイートを投稿しました",
		"tweet": gin.H{
			"id":         tweet.ID,
			"user_id":    tweet.UserID,
			"content":    tweet.Content,
			"created_at": tweet.CreatedAt,
		},
	})
}

// GetMyTweets は自分のツイート一覧を取得
func (ctrl *TweetController) GetMyTweets(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	tweets, err := ctrl.tweetService.GetUserTweets(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tweets": tweets,
	})
}

// GetAllTweets は全ツイートを取得
func (ctrl *TweetController) GetAllTweets(c *gin.Context) {
	tweets, err := ctrl.tweetService.GetAllTweets(c.Request.Context(), 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tweets": tweets,
	})
}

// DeleteTweet はツイートを削除
func (ctrl *TweetController) DeleteTweet(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// URLパラメータからツイートIDを取得
	var tweetID int32
	if _, err := c.Params.Get("id"); !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ツイートIDが必要です"})
		return
	}

	if err := ctrl.tweetService.DeleteTweet(c.Request.Context(), tweetID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ツイートを削除しました",
	})
}
