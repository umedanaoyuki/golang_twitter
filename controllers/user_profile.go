package controllers

import (
	"errors"
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserProfileController struct {
	userProfileService services.UserProfileService
}

func NewUserProfileController(userProfileService services.UserProfileService) *UserProfileController {
	return &UserProfileController{userProfileService: userProfileService}
}

type getUserProfileURI struct {
	UserID int32 `uri:"user_id" binding:"required,min=1"`
}

type CreateUserProfileBody struct {
	Name     string `json:"name" binding:"max=20" example:"たろう"`
	Bio      string `json:"bio" binding:"max=160" example:"自己紹介文です"`
	ImageURL string `json:"image_url" example:"https://example.com/avatar.png"`
	Location string `json:"location" binding:"max=10" example:"Tokyo"`
}

type UpdateUserProfileBody struct {
	Name     string `json:"name" binding:"max=20" example:"たろう"`
	Bio      string `json:"bio" binding:"max=160" example:"自己紹介文です"`
	ImageURL string `json:"image_url" example:"https://example.com/avatar.png"`
	Location string `json:"location" binding:"max=10" example:"Tokyo"`
}

type PresignProfileImageInput struct {
	ContentType string `json:"content_type" binding:"required" example:"image/jpeg"`
	Size        int64  `json:"size" binding:"required,min=1" example:"102400"`
}

type CompleteProfileImageInput struct {
	Key string `json:"key" binding:"required" example:"uploads/1_abc123.jpg"`
}

// CreateUserProfile godoc
// @Summary      プロフィール作成
// @Description  ログイン中のユーザーのプロフィールを作成する
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        input  body      CreateUserProfileBody  true  "プロフィール情報"
// @Success      201    {object}  UserProfileResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      409    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /profile [post]
func (ctrl *UserProfileController) CreateUserProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var body CreateUserProfileBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := ctrl.userProfileService.CreateUserProfile(c.Request.Context(), userID, body.Name, body.Bio, body.ImageURL, body.Location)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"profile": profile})
}

// GetUserProfileByUserID godoc
// @Summary      プロフィール取得
// @Description  指定ユーザーのプロフィールを取得する
// @Tags         profiles
// @Produce      json
// @Param        user_id  path      int  true  "ユーザーID"
// @Success      200      {object}  UserProfileResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users/{user_id}/profile [get]
func (ctrl *UserProfileController) GetUserProfileByUserID(c *gin.Context) {
	var uri getUserProfileURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := ctrl.userProfileService.GetUserProfileByUserID(c.Request.Context(), uri.UserID)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

// UpdateUserProfile godoc
// @Summary      プロフィール更新
// @Description  ログイン中のユーザーのプロフィールを更新する
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        input  body      UpdateUserProfileBody  true  "プロフィール情報"
// @Success      200    {object}  UserProfileResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      404    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /profile [put]
func (ctrl *UserProfileController) UpdateUserProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var body UpdateUserProfileBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := ctrl.userProfileService.UpdateUserProfile(c.Request.Context(), userID, body.Name, body.Bio, body.ImageURL, body.Location)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

// PresignProfileImage godoc
// @Summary      プロフィール画像アップロード許可
// @Description  認証済みユーザー向けに S3 への直接アップロード用 URL と key を発行する
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        input  body      PresignProfileImageInput  true  "画像メタ情報"
// @Success      200    {object}  PresignProfileImageResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /profile-image/presign [post]
func (ctrl *UserProfileController) PresignProfileImage(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var input PresignProfileImageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.userProfileService.PresignImageUpload(
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

// CompleteProfileImage godoc
// @Summary      プロフィール画像アップロード完了
// @Description  S3 へのアップロード完了後、key を確認してプロフィール画像を保存する
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        input  body      CompleteProfileImageInput  true  "アップロード済み画像の key"
// @Success      200    {object}  UserProfileResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      404    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /profile-image/complete [post]
func (ctrl *UserProfileController) CompleteProfileImage(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var input CompleteProfileImageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := ctrl.userProfileService.CompleteImageUpload(c.Request.Context(), userID, input.Key)
	if err != nil {
		if errors.Is(err, services.ErrUserProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
