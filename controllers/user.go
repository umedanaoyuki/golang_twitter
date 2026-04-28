package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) GetUserDetailByUserID(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	user, err := ctrl.userService.GetUserDetailByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}