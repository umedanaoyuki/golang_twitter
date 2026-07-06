package controllers

import (
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
