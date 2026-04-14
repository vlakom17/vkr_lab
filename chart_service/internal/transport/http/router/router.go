package router

import (
	"charts-chart-service/internal/client"
	"charts-chart-service/internal/transport/http/handler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	chartHandler *handler.ChartHandler,
	userClient *client.UserClient,
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

	SetupChartRoutes(r, chartHandler, userClient)

	return r
}
