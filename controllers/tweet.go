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

type CreateImageTweetInput struct {
	ImageURL string `json:"image_url" binding:"required" example:"https://example.com/image.jpg"`
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

// CreateImageTweet godoc
// @Summary      画像投稿
// @Description  認証済みユーザーとして画像ファイルを投稿する
// @Tags         tweets
// @Accept       multipart/form-data
// @Produce      json
// @Security     SessionAuth
// @Param        image  formData  file  true  "画像ファイル（JPEG/PNG・5MB以下）"
// @Success      201    {object}  CreateImageTweetResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /tweets-image [post]
func (ctrl *TweetController) CreateImageTweet(c *gin.Context) {
	// ミドルウェアからユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルを送信してください"})
		return
	}

	opened, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルを読み込めませんでした"})
		return
	}
	defer opened.Close()

	// Tweetの作成
	tweet, err := ctrl.tweetService.CreateImageTweet(
		c.Request.Context(),
		userID,
		file.Header.Get("Content-Type"),
		opened,
		file.Size,
	)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			// サービス層でバリデーション済みの日本語メッセージを返す
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tweet": gin.H{
			"user_id":    tweet.UserID,
			"image_url":  tweet.ImageUrl,
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
		if errors.Is(err, errors.New("Tweetが見つかりませんでした")) {
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