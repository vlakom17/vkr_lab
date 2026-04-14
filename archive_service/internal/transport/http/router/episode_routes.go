package router

import (
	"charts-archive-service/internal/transport/http/handler"
	"charts-archive-service/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupEpisodeRoutes(
	r *gin.Engine,
	episodeHandler *handler.EpisodeHandler,
	internalKey string,
) {

	episodes := r.Group("/episodes")
	{
		episodes.GET("", episodeHandler.GetLatestEpisodesPage)
		episodes.GET("/latest", episodeHandler.GetLatestEpisodes)
		episodes.GET("/chart/:chart_id", episodeHandler.GetEpisodesByChart)
		episodes.GET("/:id", episodeHandler.GetEpisode)
	}

	internal := r.Group("/internal/episodes")
	internal.Use(middleware.InternalAuthMiddleware(internalKey))
	{
		internal.GET("/nearest-left", episodeHandler.GetNearestLeftEpisode)
		internal.GET("/latest", episodeHandler.GetLatestEpisodesInternal)
	}
}
