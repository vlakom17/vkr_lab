package middleware

import (
	"charts-analytics-service/internal/client"
	"strings"

	"github.com/gin-gonic/gin"
)

func OptionalAuthMiddleware(userClient *client.UserClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

			if token != "" {
				userID, err := userClient.GetUserIDByToken(
					c.Request.Context(),
					token,
				)

				if err == nil {
					c.Set("user_id", userID)
				}
			}
		}

		c.Next()
	}
}
