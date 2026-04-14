package router

import (
	"charts-archive-service/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupTrackRoutes(
	r *gin.Engine,
	trackHandler *handler.TrackHandler,
) {

	tracks := r.Group("/tracks")

	tracks.GET("/search", trackHandler.SearchTracks)
}
