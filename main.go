package main

import "github.com/gin-gonic/gin"

func main() {
  router := gin.Default()
  router.GET("/health_check", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "status": "ok",
    })
  })
  router.Run() // デフォルトで0.0.0.0:8080で待機します
}