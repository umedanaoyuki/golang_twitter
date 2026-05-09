package controllers

import (
	"golang_twitter/services"

	"github.com/gin-gonic/gin"
)

type GroupController struct {
	GroupService services.GroupService
}

func NewGroupController(GroupService services.GroupService) *GroupController {
	return &GroupController{GroupService: GroupService}
}

func (ctrl *GroupController) CreateGroup(c *gin.Context) {
	return;
}

func (ctrl *GroupController) GetGroups(c *gin.Context) {
	return;
}