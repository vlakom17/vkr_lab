package router

import (
	"charts-user-service/internal/service"
	"charts-user-service/internal/transport/http/handler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authService *service.AuthService,
	userService *service.UserService,
	internalKey string,
) *gin.Engine {

	r := gin.Default()

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	SetupAuthRoutes(r, authHandler, authService, internalKey)
	SetupUserRoutes(r, userHandler, authService)

	return r
}
