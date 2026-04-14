package router

import (
	"charts-analytics-service/internal/client"
	"charts-analytics-service/internal/transport/http/handler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	reactionHandler *handler.ReactionHandler,
	recommendationHandler *handler.RecommendationHandler,
	userClient *client.UserClient,
	internalKey string,
) *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	SetupReactionRoutes(r, reactionHandler, internalKey)

	SetupRecommendationRoutes(r, recommendationHandler, userClient)

	return r
}
