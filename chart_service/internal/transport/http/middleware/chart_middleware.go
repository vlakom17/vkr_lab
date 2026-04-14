package middleware

import (
	"strings"

	"charts-chart-service/internal/client"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userClient *client.UserClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "invalid Authorization format"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if token == "" {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		userID, err := userClient.GetUserIDByToken(
			c.Request.Context(),
			token,
		)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
