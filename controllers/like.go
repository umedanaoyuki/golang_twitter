package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeController struct {
	LikeService services.LikeService
}

type LikeTweetId struct {
	TweetID int32 `uri:"id" binding:"required,min=1"`
}

func NewLikeController(LikeService services.LikeService) *LikeController {
	return &LikeController{LikeService: LikeService}
}

// CreateLike godoc
// @Summary      いいね
// @Description  指定ツイートにいいねする
// @Tags         likes
// @Produce      json
// @Security     SessionAuth
// @Param        id   path      int  true  "ツイートID"
// @Success      200  {object}  StatusOKResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tweets/{id}/like [post]
func (ctrl *LikeController) CreateLike(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri LikeTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.LikeService.CreateLike(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DeleteLike godoc
// @Summary      いいね解除
// @Description  指定ツイートのいいねを解除する
// @Tags         likes
// @Produce      json
// @Security     SessionAuth
// @Param        id   path      int  true  "ツイートID"
// @Success      200  {object}  StatusOKResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tweets/{id}/like [delete]
func (ctrl *LikeController) DeleteLike(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri LikeTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.LikeService.DeleteLike(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}