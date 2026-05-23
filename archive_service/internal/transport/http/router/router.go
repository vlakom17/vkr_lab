package router

import (
	"time"

	"charts-archive-service/internal/transport/http/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	episodeHandler *handler.EpisodeHandler,
	trackHandler *handler.TrackHandler,
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

	SetupEpisodeRoutes(r, episodeHandler, internalKey)
	SetupTrackRoutes(r, trackHandler)

	return r
}
