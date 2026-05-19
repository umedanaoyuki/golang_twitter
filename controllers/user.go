package controllers

import (
	"errors"
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

type getUserByIDURI struct {
	UserID int32 `uri:"user_id" binding:"required,min=1"`
}

// GetUserByID godoc
// @Summary      ユーザー情報取得
// @Description  指定IDのユーザー情報を取得する
// @Tags         users
// @Produce      json
// @Param        user_id  path      int  true  "ユーザーID"
// @Success      200      {object}  GetUserResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users/{user_id} [get]
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

// DeleteUser godoc
// @Summary      退会
// @Description  ログイン中のユーザーアカウントを削除する
// @Tags         users
// @Produce      json
// @Security     SessionAuth
// @Success      200  {object}  StatusOKResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /user [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	err = ctrl.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}