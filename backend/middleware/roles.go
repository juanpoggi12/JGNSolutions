package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso solo para administradores"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CheckUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso solo para usuarios"})
			c.Abort()
			return
		}
		c.Next()
	}
}
