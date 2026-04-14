package router

import (
	"charts-chart-service/internal/client"
	"charts-chart-service/internal/transport/http/handler"
	"charts-chart-service/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChartRoutes(
	r *gin.Engine,
	chartHandler *handler.ChartHandler,
	userClient *client.UserClient,
) {

	charts := r.Group("/charts")
	charts.GET("/popular", chartHandler.GetMostPopularCharts)
	charts.GET("/guest/:id", chartHandler.GetChartByIDWithoutView)
	charts.GET("/:id",
		middleware.OptionalAuthMiddleware(userClient),
		chartHandler.GetChartByID,
	)

	charts.Use(middleware.AuthMiddleware(userClient))

	charts.GET("/me", chartHandler.GetMyChart)
	charts.POST("/", chartHandler.CreateChart)
	charts.PATCH("/:id", chartHandler.PatchChart)
	charts.POST("/:id/reaction", chartHandler.SetReaction)
	charts.POST("/:id/episode", chartHandler.CreateEpisode)
	charts.GET("/me/likes", chartHandler.GetMyLikedCharts)
	charts.GET("/me/dislikes", chartHandler.GetMyDislikedCharts)

	internal := r.Group("/internal/charts")

	internal.POST("/by-ids", chartHandler.GetChartsByIDs)
	internal.POST("/ids-by-genre", chartHandler.GetChartIDsByGenre)
	internal.POST("/genres-by-ids", chartHandler.GetGenresByChartIDs)
}
