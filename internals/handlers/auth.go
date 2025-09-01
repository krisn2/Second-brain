package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "register endpoint",
	})
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "login endpoint",
	})
}
