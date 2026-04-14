package main

import (
	"log"

	"charts-chart-service/internal/client"
	"charts-chart-service/internal/config"
	"charts-chart-service/internal/infrastructure/db"
	"charts-chart-service/internal/infrastructure/mb"
	"charts-chart-service/internal/producers/kafka"
	"charts-chart-service/internal/repository/postgres"
	"charts-chart-service/internal/service"
	"charts-chart-service/internal/transport/http/handler"
	"charts-chart-service/internal/transport/http/router"
)

func main() {

	cfg := config.Load()

	pg := db.NewPostgresPool(cfg.DatabaseURL)
	db.RunMigrations(cfg.DatabaseURL)

	kafkaConn := mb.NewKafka(cfg.KafkaBroker)
	defer kafkaConn.Close()

	producer := kafka.NewKafkaProducer(kafkaConn)
	userClient := client.NewUserClient(cfg.UserServiceURL, cfg.InternalAPIKey)
	analyticsClient := client.NewAnalyticsClient(cfg.AnalyticsServiceURL, cfg.InternalAPIKey)

	repo := postgres.NewChartRepository(pg)

	chartService := service.NewChartService(repo, userClient, analyticsClient, producer)

	chartHandler := handler.NewChartHandler(chartService)

	r := router.SetupRouter(chartHandler, userClient)
	r.SetTrustedProxies(nil)

	log.Println("chart-service started on :8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
