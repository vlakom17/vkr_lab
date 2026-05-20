package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"charts-analytics-service/internal/client"
	"charts-analytics-service/internal/config"
	kafkaConsumer "charts-analytics-service/internal/consumers/kafka"
	"charts-analytics-service/internal/infrastructure/db"
	"charts-analytics-service/internal/infrastructure/mb"
	"charts-analytics-service/internal/repository/postgres"
	"charts-analytics-service/internal/service"
	"charts-analytics-service/internal/transport/http/handler"
	"charts-analytics-service/internal/transport/http/router"
)

func main() {
	cfg := config.Load()

	pg := db.NewPostgresPool(cfg.DatabaseURL)
	db.RunMigrations(cfg.DatabaseURL)

	reactionRepo := postgres.NewReactionRepository(pg)

	archiveClient := client.NewArchiveClient(cfg.ArchiveServiceURL, cfg.InternalAPIKey)
	userClient := client.NewUserClient(cfg.UserServiceURL, cfg.InternalAPIKey)

	reactionService := service.NewReactionService(reactionRepo, userClient)
	recommendationService := service.NewRecommendationService(reactionRepo, archiveClient)

	consumer := mb.NewConsumer(cfg.KafkaBroker, mb.KafkaTopic, cfg.KafkaGroupID)
	kafkaCons := kafkaConsumer.NewKafkaConsumer(consumer)

	reactionHandler := handler.NewReactionHandler(reactionService)
	recommendationHandler := handler.NewRecommendationHandler(recommendationService)

	r := router.SetupRouter(
		reactionHandler,
		recommendationHandler,
		userClient,
		cfg.InternalAPIKey,
	)

	r.SetTrustedProxies(nil)

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go kafkaCons.ConsumeReactions(ctx, reactionService.HandleReactionEvent)

	go func() {
		log.Println("analytics-service started on :8080")
		if err := r.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	sig := <-sigChan
	log.Println("shutdown signal:", sig)

	cancel()

	time.Sleep(2 * time.Second)

	if err := consumer.Close(); err != nil {
		log.Println("error closing kafka consumer:", err)
	}

	log.Println("service stopped")
}
