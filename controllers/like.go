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

// GetUserLikes godoc
// @Summary      ユーザーのいいね一覧取得
// @Description  指定ユーザーがいいねしたツイートをカーソルページネーションで取得する
// @Tags         likes
// @Produce      json
// @Param        user_id  path      int   true   "ユーザーID"
// @Param        cursor   query     int   false  "ページネーションカーソル（最後に取得したいいねID）"
// @Param        limit    query     int   false  "取得件数（1〜100、デフォルト20）"  default(20)
// @Success      200      {object}  GetUserLikesResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users/{user_id}/likes [get]
func (ctrl *LikeController) GetUserLikes(c *gin.Context) {
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

	// いいね一覧（ツイート詳細付き）を取得
	likes, err := ctrl.LikeService.GetUserLikesWithCursor(
		c.Request.Context(),
		uriParams.UserID,
		queryParams.Cursor,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 次のカーソルを計算（likes テーブルの id）
	var nextCursor *int32
	if len(likes) > 0 {
		lastID := likes[len(likes)-1].LikeRowID
		nextCursor = &lastID
	}

	c.JSON(http.StatusOK, gin.H{
		"likes":       likes,
		"next_cursor": nextCursor,
		"has_more":    len(likes) == int(limit),
	})
}