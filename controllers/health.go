package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary      ヘルスチェック
// @Description  サーバーの稼働状態を確認する
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthCheckResponse
// @Router       /health_check [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
