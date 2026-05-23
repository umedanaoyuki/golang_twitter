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
	Content string `json:"content" binding:"required" example:"メッセージ本文"`
}

type CreateMessageURI struct {
	GroupID int32 `uri:"group_id" binding:"required,min=1"`
}

// CreateMessage godoc
// @Summary      メッセージ送信
// @Description  グループ内にメッセージを送信する（グループメンバーのみ）
// @Tags         messages
// @Accept       json
// @Produce      json
// @Security     SessionAuth
// @Param        group_id  path      int                 true  "グループID"
// @Param        input     body      CreateMessageInput  true  "メッセージ内容"
// @Success      201       {object}  CreateMessageResponse
// @Failure      400       {object}  ErrorResponse
// @Failure      401       {object}  ErrorResponse
// @Failure      403       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Router       /groups/{group_id}/messages [post]
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
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
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

// GetMessages godoc
// @Summary      メッセージ一覧取得
// @Description  グループ内のメッセージ一覧を取得する（グループメンバーのみ）
// @Tags         messages
// @Produce      json
// @Security     SessionAuth
// @Param        group_id  path      int  true  "グループID"
// @Success      200       {object}  GetMessagesResponse
// @Failure      400       {object}  ErrorResponse
// @Failure      401       {object}  ErrorResponse
// @Failure      403       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Router       /groups/{group_id}/messages [get]
func (ctrl *MessageController) GetMessages(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var uri CreateMessageURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なパスパラメータです"})
		return
	}

	messages, err := ctrl.messageService.GetMessages(c.Request.Context(), userID, uri.GroupID)
	if err != nil {
		if _, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}