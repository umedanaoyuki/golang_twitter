package controllers

import (
	"errors"
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

type getUserByIDURI struct {
	UserID int32 `uri:"user_id" binding:"required,min=1"`
}

// GetUserByID は /users/:user_id で指定されたユーザーの公開情報を返す。
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	var uri getUserByIDURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userService.GetUserDetailByUserID(c.Request.Context(), uri.UserID)
	if err != nil {
		if errors.Is(err, errors.New(err.Error())){
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}