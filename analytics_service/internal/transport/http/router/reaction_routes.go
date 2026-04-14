package router

import (
	"charts-analytics-service/internal/transport/http/handler"
	"charts-analytics-service/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupReactionRoutes(
	r *gin.Engine,
	h *handler.ReactionHandler,
	internalKey string,
) {

	public := r.Group("/analysis")
	{
		public.GET("/:chart_id", h.GetReactionStats)
	}

	internal := r.Group("/internal/analysis")
	internal.Use(middleware.InternalAuthMiddleware(internalKey))
	{
		internal.GET("/charts/popular", h.GetMostPopularChartIDs)
		internal.GET("/users/:user_id/likes", h.GetUserLikedChartIDs)
		internal.GET("/users/:user_id/dislikes", h.GetUserDislikedChartIDs)
	}
}
