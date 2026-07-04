package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	CommentService services.CommentService
}

type CommentTweetId struct {
	TweetID int32 `uri:"id" binding:"required,min=1"`
}

type GetCommentsQuery struct {
	Cursor *int32 `form:"cursor" binding:"omitempty,min=1"`
	Limit  int32  `form:"limit" binding:"omitempty,min=1,max=100"`
}

func NewCommentController(CommentService services.CommentService) *CommentController {
	return &CommentController{CommentService: CommentService}
}


func (ctrl *CommentController) CreateComment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri CommentTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.CommentService.CreateComment(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}


func (ctrl *CommentController) DeleteComment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri CommentTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.CommentService.DeleteComment(c.Request.Context(), userID, uri.TweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetComments godoc
// @Summary      コメント一覧取得
// @Description  指定ツイートのコメントをカーソルページネーションで取得する
// @Tags         comments
// @Produce      json
// @Param        id      path      int   true   "ツイートID"
// @Param        cursor  query     int   false  "ページネーションカーソル（最後に取得したコメントID）"
// @Param        limit   query     int   false  "取得件数（1〜100、デフォルト20）"  default(20)
// @Success      200     {object}  GetCommentsResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /tweets/{id}/comments [get]
func (ctrl *CommentController) GetComments(c *gin.Context) {
	var uri CommentTweetId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なパスパラメータです"})
		return
	}

	var queryParams GetCommentsQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なクエリパラメータです"})
		return
	}

	limit := queryParams.Limit
	if limit == 0 {
		limit = 20
	}

	comments, err := ctrl.CommentService.GetCommentsByTweetIDWithCursor(
		c.Request.Context(),
		uri.TweetID,
		queryParams.Cursor,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var nextCursor *int32
	if len(comments) > 0 {
		lastID := comments[len(comments)-1].ID
		nextCursor = &lastID
	}

	c.JSON(http.StatusOK, gin.H{
		"comments":    comments,
		"next_cursor": nextCursor,
		"has_more":    len(comments) == int(limit),
	})
}