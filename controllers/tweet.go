package controllers

import (
	"net/http"
	"strconv"

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
	// URLパラメータからuser_idを取得
	userIDStr := c.Param("user_id")
	userID64, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーIDです"})
		return
	}
	userID := int32(userID64)

	// クエリパラメータからcursorとlimitを取得
	var cursor *int32
	if cursorStr := c.Query("cursor"); cursorStr != "" {
		cursor64, err := strconv.ParseInt(cursorStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なカーソルです"})
			return
		}
		cursorVal := int32(cursor64)
		cursor = &cursorVal
	}

	// デフォルトで20件取得
	limit := int32(20)
	if limitStr := c.Query("limit"); limitStr != "" {
		limit64, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なlimitです"})
			return
		}
		limit = int32(limit64)
	}

	// ツイート一覧を取得
	tweets, err := ctrl.tweetService.GetUserTweetsWithCursor(c.Request.Context(), userID, cursor, limit)
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