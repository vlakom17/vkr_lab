package router

import (
	"charts-user-service/internal/service"
	"charts-user-service/internal/transport/http/handler"
	"charts-user-service/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	authService *service.AuthService,
) {

	users := r.Group("/users")

	users.GET("/:id", userHandler.GetUserProfile)

	users.Use(middleware.AuthMiddleware(authService))

	users.GET("/me", userHandler.GetMyProfile)
	users.PATCH("/me", userHandler.UpdateUser)
}
