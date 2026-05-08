package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RetweetController struct {
	RetweetService services.RetweetService
}

type RetweetTweetId struct {
	TweetID int32 `uri:"id" binding:"required,min=1"`
}

type GetUserRetweetsByUserId struct {
	UserID int32 `uri:"user_id" binding:"required,min=1"`
}

func NewRetweetController(RetweetService services.RetweetService) *RetweetController {
	return &RetweetController{RetweetService: RetweetService}
}

func (ctrl *RetweetController) CreateRetweet(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri RetweetTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.RetweetService.CreateRetweet(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (ctrl *RetweetController) DeleteRetweet(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri RetweetTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.RetweetService.DeleteRetweet(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (ctrl *RetweetController) GetUserRetweets(c *gin.Context) {
	// パスパラメータのバリデーション
	var uriParams GetUserTweetsUri
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なパスパラメータです"})
		return
	}

	// クエリパラメータのバリデーション（デフォルト値を設定）
	var queryParams GetUserTweetsQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なクエリパラメータです"})
		return
	}

	// limitのデフォルト値を設定
	limit := queryParams.Limit
	if limit == 0 {
		limit = 20
	}

	// リツイート一覧を取得
	tweets, err := ctrl.RetweetService.GetUserRetweetsWithCursor(
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