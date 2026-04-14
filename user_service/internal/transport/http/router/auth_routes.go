package router

import (
	"charts-user-service/internal/service"
	"charts-user-service/internal/transport/http/handler"
	"charts-user-service/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
	internalKey string,
) {
	auth := r.Group("/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	auth.Use(middleware.AuthMiddleware(authService))
	auth.POST("/logout", authHandler.Logout)

	internal := r.Group("/internal")
	internal.Use(middleware.InternalAuthMiddleware(internalKey))
	internal.Use(middleware.AuthMiddleware(authService))
	internal.GET("/auth/user", authHandler.GetUserIDByToken)
}
