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