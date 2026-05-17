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
	Name          string  `json:"name" binding:"required"`
	MemberUserIDs []int32 `json:"member_user_ids"`
}

func NewGroupController(GroupService services.GroupService) *GroupController {
	return &GroupController{GroupService: GroupService}
}

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
