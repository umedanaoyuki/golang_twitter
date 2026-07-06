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

type DeleteCommentURI struct {
	TweetID   int32 `uri:"id" binding:"required,min=1"`
	CommentID int32 `uri:"comment_id" binding:"required,min=1"`
}

type CreateCommentBody struct {
	Content string `json:"content" binding:"required" example:"いいツイートですね！"`
}

type GetCommentsQuery struct {
	Cursor *int32 `form:"cursor" binding:"omitempty,min=1"`
	Limit  int32  `form:"limit" binding:"omitempty,min=1,max=100"`
}

func NewCommentController(CommentService services.CommentService) *CommentController {
	return &CommentController{CommentService: CommentService}
}

// CreateComment godoc
// @Summary      コメント作成
// @Description  指定ツイートにコメントする
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        id     path      int                 true  "ツイートID"
// @Param        input  body      CreateCommentBody   true  "コメント内容"
// @Success      201    {object}  CreateCommentResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /tweets/{id}/comments [post]
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

	var body CreateCommentBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := ctrl.CommentService.CreateComment(c.Request.Context(), userID, uri.TweetID, body.Content)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}

// DeleteComment godoc
// @Summary      コメント削除
// @Description  指定ツイートの自分のコメントを削除する
// @Tags         comments
// @Produce      json
// @Security     SessionAuth
// @Param        id          path      int  true  "ツイートID"
// @Param        comment_id  path      int  true  "コメントID"
// @Success      200         {object}  StatusOKResponse
// @Failure      400         {object}  ErrorResponse
// @Failure      401         {object}  ErrorResponse
// @Failure      404         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /tweets/{id}/comments/{comment_id} [delete]
func (ctrl *CommentController) DeleteComment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri DeleteCommentURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.CommentService.DeleteComment(c.Request.Context(), userID, uri.TweetID, uri.CommentID); err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
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
