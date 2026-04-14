package router

import (
	"charts-analytics-service/internal/client"
	"charts-analytics-service/internal/transport/http/handler"
	"charts-analytics-service/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRecommendationRoutes(
	r *gin.Engine,
	h *handler.RecommendationHandler,
	userClient *client.UserClient,
) {

	group := r.Group("/analysis")

	group.GET(
		"/recommendations",
		middleware.OptionalAuthMiddleware(userClient),
		h.GetRecommendations,
	)
}
