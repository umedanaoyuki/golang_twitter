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

// CreateRetweet godoc
// @Summary      リツイート
// @Description  指定ツイートをリツイートする
// @Tags         retweets
// @Produce      json
// @Security     SessionAuth
// @Param        id   path      int  true  "ツイートID"
// @Success      200  {object}  StatusOKResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tweets/{id}/retweet [post]
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

// DeleteRetweet godoc
// @Summary      リツイート解除
// @Description  指定ツイートのリツイートを解除する
// @Tags         retweets
// @Produce      json
// @Security     SessionAuth
// @Param        id   path      int  true  "ツイートID"
// @Success      200  {object}  StatusOKResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tweets/{id}/retweet [delete]
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

// GetUserRetweets godoc
// @Summary      ユーザーのリツイート一覧取得
// @Description  指定ユーザーのリツイートをカーソルページネーションで取得する
// @Tags         retweets
// @Produce      json
// @Param        user_id  path      int   true   "ユーザーID"
// @Param        cursor   query     int   false  "ページネーションカーソル（最後に取得したリツイートID）"
// @Param        limit    query     int   false  "取得件数（1〜100、デフォルト20）"  default(20)
// @Success      200      {object}  GetUserRetweetsResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users/{user_id}/retweets [get]
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

	// リツイート一覧（ツイート詳細付き）を取得
	retweets, err := ctrl.RetweetService.GetUserRetweetsWithCursor(
		c.Request.Context(),
		uriParams.UserID,
		queryParams.Cursor,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 次のカーソルを計算（retweets テーブルの id）
	var nextCursor *int32
	if len(retweets) > 0 {
		lastID := retweets[len(retweets)-1].RetweetRowID
		nextCursor = &lastID
	}

	c.JSON(http.StatusOK, gin.H{
		"retweets":    retweets,
		"next_cursor": nextCursor,
		"has_more":    len(retweets) == int(limit),
	})
}