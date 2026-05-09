package controllers

import (
	"golang_twitter/services"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	MessageService services.MessageService
}

func NewMessageController(MessageService services.MessageService) *MessageController {
	return &MessageController{MessageService: MessageService}
}

func (ctrl *MessageController) CreateMessage(c *gin.Context) {
	return;
}

func (ctrl *MessageController) GetMessages(c *gin.Context) {
	return;
}