package controllers

import (
	"log"
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

	log.Println(tweet.UserID)

	c.JSON(http.StatusCreated, gin.H{
		"tweet": gin.H{
			"user_id":    tweet.UserID,
			"content":    tweet.Content,
		},
	})
}


func (ctrl *TweetController) GetTweetByID(c *gin.Context) {
	id := c.Param("id")

	// strconv.ParseInt(文字列, 進数, ビットサイズ)
	// bitSize=32 で int32 を指定
	i64, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDが不正です"})
		return
	}
	
	// int64型で返るため、int32にキャスト
	i32 := int32(i64)

	tweet, err := ctrl.tweetService.GetTweetByID(c.Request.Context(), int32(i32))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tweetが見つかりませんでした"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tweet": tweet,
	})
}