package middleware

import (
	"net/http"
	"strings"

	"charts-user-service/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := authService.GetUserIDByToken(
			c.Request.Context(),
			token,
		)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
