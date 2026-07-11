package controllers

import (
	"errors"
	"net/http"

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
	Content string `json:"content" binding:"required" example:"Hello, world!"`
}

type DeleteTweetInput struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}
type CreateImageTweetInput struct {
	ImageURL string `json:"image_url" binding:"required" example:"https://example.com/image.jpg"`
}

type PresignImageTweetInput struct {
	ContentType string `json:"content_type" binding:"required" example:"image/jpeg"`
	Size        int64  `json:"size" binding:"required,min=1" example:"102400"`
}

type CompleteImageTweetInput struct {
	Key string `json:"key" binding:"required" example:"uploads/1_abc123.jpg"`
}

type GetTweetByIdInput struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}
type GetUserTweetsUri struct {
	UserID int32 `uri:"user_id" binding:"required,min=1"`
}

type GetUserTweetsQuery struct {
	Cursor *int32 `form:"cursor" binding:"omitempty,min=1"`
	Limit  int32  `form:"limit" binding:"omitempty,min=1,max=100"`
}

// CreateTweet godoc
// @Summary      ツイート投稿
// @Description  認証済みユーザーとしてツイートを投稿する
// @Tags         tweets
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        input  body      CreateTweetInput  true  "ツイート内容"
// @Success      201    {object}  CreateTweetResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /tweets [post]
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

	c.JSON(http.StatusCreated, gin.H{
		"tweet": gin.H{
			"user_id":    tweet.UserID,
			"content":    tweet.Content,
		},
	})
}


// DeleteTweet godoc
// @Summary      ツイート削除
// @Description  認証済みユーザーが自分のツイートを削除する（他人のツイートは削除できない）
// @Tags         tweets
// @Produce      json
// @Security     SessionAuth
// @Param        id   path      int  true  "ツイートID"
// @Success      200  {object}  StatusOKResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tweets/{id} [delete]
func (ctrl *TweetController) DeleteTweet(c *gin.Context) {
	var input DeleteTweetInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	if err := ctrl.tweetService.DeleteTweet(c.Request.Context(), userID, input.ID); err != nil {
		if errors.Is(err, services.ErrTweetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// PresignImageTweet godoc
// @Summary      画像アップロード許可
// @Description  認証済みユーザー向けに S3 への直接アップロード用 URL と key を発行する
// @Tags         tweets
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        input  body      PresignImageTweetInput  true  "画像メタ情報"
// @Success      200    {object}  PresignImageTweetResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /tweets-image/presign [post]
func (ctrl *TweetController) PresignImageTweet(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var input PresignImageTweetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.tweetService.PresignImageUpload(
		c.Request.Context(),
		userID,
		input.ContentType,
		input.Size,
	)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":        result.Key,
		"upload_url": result.UploadURL,
		"public_url": result.PublicURL,
	})
}

// CompleteImageTweet godoc
// @Summary      画像アップロード完了
// @Description  S3 へのアップロード完了後、key を確認して画像ツイートを作成する
// @Tags         tweets
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        input  body      CompleteImageTweetInput  true  "アップロード済み画像の key"
// @Success      201    {object}  CreateImageTweetResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /tweets-image/complete [post]
func (ctrl *TweetController) CompleteImageTweet(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var input CompleteImageTweetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet, err := ctrl.tweetService.CompleteImageUpload(c.Request.Context(), userID, input.Key)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tweet": gin.H{
			"user_id":   tweet.UserID,
			"image_url": tweet.ImageUrl,
		},
	})
}

// GetUserTweets godoc
// @Summary      ユーザーのツイート一覧取得
// @Description  指定ユーザーのツイートをカーソルページネーションで取得する
// @Tags         tweets
// @Produce      json
// @Param        user_id  path      int   true   "ユーザーID"
// @Param        cursor   query     int   false  "ページネーションカーソル（最後に取得したツイートID）"
// @Param        limit    query     int   false  "取得件数（1〜100、デフォルト20）"  default(20)
// @Success      200      {object}  GetUserTweetsResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users/{user_id}/tweets [get]
func (ctrl *TweetController) GetUserTweets(c *gin.Context) {
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

	// ツイート一覧を取得
	tweets, err := ctrl.tweetService.GetUserTweetsWithCursor(
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

// GetTweetByID godoc
// @Summary      ツイート取得
// @Description  指定IDのツイートを取得する
// @Tags         tweets
// @Produce      json
// @Param        id   path      int  true  "ツイートID"
// @Success      200  {object}  GetTweetResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tweets/{id} [get]
func (ctrl *TweetController) GetTweetByID(c *gin.Context) {
	var input GetTweetByIdInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なパスパラメータです"})
		return
	}

	tweet, err := ctrl.tweetService.GetTweetByID(c.Request.Context(), input.Id)
	if err != nil {
		if errors.Is(err, services.ErrTweetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tweet": tweet,
	})
}