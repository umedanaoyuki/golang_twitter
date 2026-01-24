package main

import (
	"golang_twitter/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  router.GET("/health_check", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "status": "ok",
    })
  })

router.POST("/register", controllers.Register)

//   デフォルト 0.0.0.0:8080
  router.Run()
}