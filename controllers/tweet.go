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
		"tweet": tweetResponseFromTweet(tweet),
	})
}

// CreateTweetWithImage godoc
// @Summary      画像付きツイート投稿
// @Description  認証済みユーザーとして画像付きツイートを投稿する（本文は任意）
// @Tags         tweets
// @Accept       multipart/form-data
// @Produce      json
// @Security     SessionAuth
// @Param        content  formData  string  false  "ツイート本文（任意・140文字以内）"
// @Param        image    formData  file    true   "画像ファイル（JPEG/PNG/GIF/WebP・5MB以下）"
// @Success      201      {object}  CreateTweetResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /tweets/with-image [post]
func (ctrl *TweetController) CreateTweetWithImage(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	content := c.PostForm("content")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルを指定してください"})
		return
	}

	opened, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルを読み込めませんでした"})
		return
	}
	defer opened.Close()

	tweet, err := ctrl.tweetService.CreateTweetWithImage(
		c.Request.Context(),
		userID,
		content,
		file.Filename,
		file.Header.Get("Content-Type"),
		opened,
		file.Size,
	)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tweet": tweetResponseFromTweet(tweet),
	})
}

func tweetResponseFromTweet(tweet *services.Tweet) gin.H {
	resp := gin.H{
		"id":       tweet.ID,
		"user_id":  tweet.UserID,
		"content":  tweet.Content,
	}
	if tweet.ImageURL != nil {
		resp["image_url"] = *tweet.ImageURL
	}
	return resp
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
	var uriParams GetUserTweetsUri
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なパスパラメータです"})
		return
	}

	var queryParams GetUserTweetsQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なクエリパラメータです"})
		return
	}

	limit := queryParams.Limit
	if limit == 0 {
		limit = 20
	}

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
