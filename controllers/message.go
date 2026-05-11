package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	messageService services.MessageService
}

func NewMessageController(MessageService services.MessageService) *MessageController {
	return &MessageController{messageService: MessageService}
}

type CreateMessageInput struct {
	Content string `json:"content" binding:"required"`
}

type CreateMessageURI struct {
	GroupID int32 `uri:"group_id" binding:"required,min=1"`
}

// メッセージ作成
func (ctrl *MessageController) CreateMessage(c *gin.Context) {
	// ミドルウェアからユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var uri CreateMessageURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なパスパラメータです"})
		return
	}

	var input CreateMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// メッセージの作成（投稿先グループはパス :group_id から取得）
	message, err := ctrl.messageService.CreateMessage(c.Request.Context(), userID, uri.GroupID, input.Content)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": gin.H{
			"user_id":    message.UserID,
			"group_id":   message.GroupID,
			"content":    message.Content,
		},
	})
}

func (ctrl *MessageController) GetMessages(c *gin.Context) {
	// パスパラメータのバリデーション（:group_id）
	var uri CreateMessageURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なパスパラメータです"})
		return
	}

	// メッセージ一覧を取得
	messages, err := ctrl.messageService.GetMessages(c.Request.Context(), uri.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}