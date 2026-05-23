package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FollowController struct {
	followService services.FollowService
}

type FollowUserId struct {
	FollowedUserID int32 `uri:"user_id" binding:"required,min=1"`
}

func NewFollowController(followService services.FollowService) *FollowController {
	return &FollowController{followService: followService}
}

// CreateFollow godoc
// @Summary      フォロー
// @Description  指定ユーザーをフォローする
// @Tags         follows
// @Produce      json
// @Security     SessionAuth
// @Param        user_id  path      int  true  "フォロー対象のユーザーID"
// @Success      200      {object}  StatusOKResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users/{user_id}/follow [post]
func (ctrl *FollowController) CreateFollow(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri FollowUserId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.followService.CreateFollow(c.Request.Context(), userID, uri.FollowedUserID); err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DeleteFollow godoc
// @Summary      フォロー解除
// @Description  指定ユーザーのフォローを解除する
// @Tags         follows
// @Produce      json
// @Security     SessionAuth
// @Param        user_id  path      int  true  "フォロー解除対象のユーザーID"
// @Success      200      {object}  StatusOKResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users/{user_id}/follow [delete]
func (ctrl *FollowController) DeleteFollow(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri FollowUserId
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.followService.DeleteFollow(c.Request.Context(), userID, uri.FollowedUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetFollowersByUserId godoc
// @Summary      フォロワー一覧取得
// @Description  指定ユーザーのフォロワーをカーソルページネーションで取得する
// @Tags         follows
// @Produce      json
// @Security     SessionAuth
// @Param        user_id  path      int   true   "ユーザーID"
// @Param        cursor   query     int   false  "ページネーションカーソル"
// @Param        limit    query     int   false  "取得件数（1〜100、デフォルト20）"  default(20)
// @Success      200      {object}  GetFollowersResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users/{user_id}/followers [get]
func (ctrl *FollowController) GetFollowersByUserId(c *gin.Context) {
	if _, err := middleware.GetUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

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

	followers, err := ctrl.followService.GetFollowersByUserIdWithCursor(
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
	if len(followers) > 0 {
		lastID := followers[len(followers)-1].ID
		nextCursor = &lastID
	}

	c.JSON(http.StatusOK, gin.H{
		"followers":   followers,
		"next_cursor": nextCursor,
		"has_more":    len(followers) == int(limit),
	})
}

// GetFollowingByUserId godoc
// @Summary      フォロー中一覧取得
// @Description  指定ユーザーがフォローしているユーザーをカーソルページネーションで取得する
// @Tags         follows
// @Produce      json
// @Security     SessionAuth
// @Param        user_id  path      int   true   "ユーザーID"
// @Param        cursor   query     int   false  "ページネーションカーソル"
// @Param        limit    query     int   false  "取得件数（1〜100、デフォルト20）"  default(20)
// @Success      200      {object}  GetFollowingResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users/{user_id}/following [get]
func (ctrl *FollowController) GetFollowingByUserId(c *gin.Context) {
	if _, err := middleware.GetUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

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

	following, err := ctrl.followService.GetFollowingByUserIdWithCursor(
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
	if len(following) > 0 {
		lastID := following[len(following)-1].ID
		nextCursor = &lastID
	}

	c.JSON(http.StatusOK, gin.H{
		"following":   following,
		"next_cursor": nextCursor,
		"has_more":    len(following) == int(limit),
	})
}