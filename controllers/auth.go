package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "data": "register endpointです",
    })
}
