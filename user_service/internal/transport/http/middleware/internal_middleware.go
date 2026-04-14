package middleware

import "github.com/gin-gonic/gin"

func InternalAuthMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Internal-Key")

		if key == "" || key != expectedKey {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Next()
	}
}
