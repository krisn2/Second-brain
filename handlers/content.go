package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddContent(c *gin.Context) {
	c.JSON(
		http.StatusOK, gin.H{
			"msg": "add content endpoint",
		},
	)
}

func GetContent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "get content endpoint",
	})
}
