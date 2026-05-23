package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GroupController struct {
	GroupService services.GroupService
}

type CreateGroupBody struct {
	Name          string  `json:"name" binding:"required" example:"グループ名"`
	MemberUserIDs []int32 `json:"member_user_ids" example:"2,3"`
}

func NewGroupController(GroupService services.GroupService) *GroupController {
	return &GroupController{GroupService: GroupService}
}

// CreateGroup godoc
// @Summary      グループ作成
// @Description  グループを作成し、メンバーを追加する
// @Tags         groups
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        input  body      CreateGroupBody  true  "グループ情報"
// @Success      201    {object}  CreateGroupResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      409    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /groups [post]
func (ctrl *GroupController) CreateGroup(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var body CreateGroupBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := ctrl.GroupService.CreateGroup(c.Request.Context(), userID, body.Name, body.MemberUserIDs)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"group": gin.H{
			"id":         group.ID,
			"name":       group.Name,
			"user_id":    group.UserID,
			"created_at": group.CreatedAt,
		},
	})
}

// GetGroups godoc
// @Summary      グループ一覧取得
// @Description  ログイン中のユーザーが所属するグループ一覧を取得する
// @Tags         groups
// @Produce      json
// @Security     SessionAuth
// @Success      200  {array}   SwaggerGroup
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /groups [get]
func (ctrl *GroupController) GetGroups(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	groups, err := ctrl.GroupService.GetGroups(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}
