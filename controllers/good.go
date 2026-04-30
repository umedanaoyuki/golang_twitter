package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GoodController struct {
	goodService services.GoodService
}

type GoodTweetId struct {
	TweetID int32 `uri:"tweet_id" binding:"required,min=1"`
}

func NewGoodController(goodService services.GoodService) *GoodController {
	return &GoodController{goodService: goodService}
}

func (ctrl *GoodController) CreateGood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var uri GoodTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.goodService.CreateGood(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (ctrl *GoodController) DeleteGood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var uri GoodTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.goodService.DeleteGood(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}