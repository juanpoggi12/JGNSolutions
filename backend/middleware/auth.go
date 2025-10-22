package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"
)

func AuthMiddleware(cfg utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "falta token"})
			c.Abort()
			return
		}

		// ✅ Usamos SplitN o Fields para tolerar espacios extra
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ParseAccessToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			c.Abort()
			return
		}

		// Guardar en el contexto (acceso global para handlers y services)
		c.Set("userId", claims.Subject)
		c.Set("role", claims.Role)

		c.Next()
	}
}

