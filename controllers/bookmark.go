package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookmarkController struct {
	BookmarkService services.BookmarkService
}

type BookmarkTweetId struct {
	TweetID int32 `uri:"id" binding:"required,min=1"`
}

func NewBookmarkController(BookmarkService services.BookmarkService) *BookmarkController {
	return &BookmarkController{BookmarkService: BookmarkService}
}

// CreateBookmark godoc
// @Summary      ブックマーク追加
// @Description  指定ツイートをブックマークする
// @Tags         bookmarks
// @Produce      json
// @Security     SessionAuth
// @Param        id   path      int  true  "ツイートID"
// @Success      200  {object}  StatusOKResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tweets/{id}/bookmark [post]
func (ctrl *BookmarkController) CreateBookmark(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri BookmarkTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.BookmarkService.CreateBookmark(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DeleteBookmark godoc
// @Summary      ブックマーク解除
// @Description  指定ツイートのブックマークを解除する
// @Tags         bookmarks
// @Produce      json
// @Security     SessionAuth
// @Param        id   path      int  true  "ツイートID"
// @Success      200  {object}  StatusOKResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tweets/{id}/bookmark [delete]
func (ctrl *BookmarkController) DeleteBookmark(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri BookmarkTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.BookmarkService.DeleteBookmark(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetBookmarksByUserId godoc
// @Summary      ブックマーク一覧取得
// @Description  ログイン中のユーザーのブックマーク一覧を取得する
// @Tags         bookmarks
// @Produce      json
// @Security     SessionAuth
// @Success      200  {array}   SwaggerBookmark
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tweets/bookmarks [get]
func (ctrl *BookmarkController) GetBookmarksByUserId(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	bookmarks, err := ctrl.BookmarkService.GetBookmarksByUserId(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookmarks)
}