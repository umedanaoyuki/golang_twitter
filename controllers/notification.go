package controllers

import (
	"golang_twitter/middleware"
	"golang_twitter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationService services.NotificationService
}

type GetNotificationsQuery struct {
	Cursor *int32 `form:"cursor" binding:"omitempty,min=1"`
	Limit  int32  `form:"limit" binding:"omitempty,min=1,max=100"`
}

func NewNotificationController(NotificationService services.NotificationService) *NotificationController {
	return &NotificationController{NotificationService: NotificationService}
}

// GetNotifications godoc
// @Summary      通知一覧取得
// @Description  ログインユーザー宛の通知（いいね・フォロー・コメント）をカーソルページネーションで取得する
// @Tags         notifications
// @Produce      json
// @Security     SessionAuth
// @Param        cursor  query     int   false  "ページネーションカーソル（最後に取得した通知ID）"
// @Param        limit   query     int   false  "取得件数（1〜100、デフォルト20）"  default(20)
// @Success      200     {object}  GetNotificationsResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      401     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /notifications [get]
func (ctrl *NotificationController) GetNotifications(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	var queryParams GetNotificationsQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なクエリパラメータです"})
		return
	}

	limit := queryParams.Limit
	if limit == 0 {
		limit = 20
	}

	notifications, err := ctrl.NotificationService.GetNotificationsWithCursor(
		c.Request.Context(),
		userID,
		queryParams.Cursor,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var nextCursor *int32
	if len(notifications) > 0 {
		lastID := notifications[len(notifications)-1].ID
		nextCursor = &lastID
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"next_cursor":   nextCursor,
		"has_more":      len(notifications) == int(limit),
	})
}

// GetNotificationCount godoc
// @Summary      未読通知件数取得
// @Description  ログインユーザー宛の未読通知の件数のみを返却する（バッジ表示などの軽量用途向け）
// @Tags         notifications
// @Produce      json
// @Security     SessionAuth
// @Success      200  {object}  GetNotificationCountResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /notifications/count [get]
func (ctrl *NotificationController) GetNotificationCount(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	count, err := ctrl.NotificationService.GetUnreadNotificationCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unread_count": count,
	})
}
