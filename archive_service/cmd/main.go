package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"charts-archive-service/internal/config"
	kafkaConsumer "charts-archive-service/internal/consumers/kafka"
	"charts-archive-service/internal/infrastructure/db"
	"charts-archive-service/internal/infrastructure/mb"
	"charts-archive-service/internal/repository/postgres"
	"charts-archive-service/internal/service"
	"charts-archive-service/internal/transport/http/handler"
	"charts-archive-service/internal/transport/http/router"
)

func main() {
	cfg := config.Load()

	pg := db.NewPostgresPool(cfg.DatabaseURL)
	db.RunMigrations(cfg.DatabaseURL)

	episodeRepo := postgres.NewEpisodeRepository(pg)
	trackRepo := postgres.NewTrackRepository(pg)

	episodeService := service.NewEpisodeService(episodeRepo, trackRepo)
	trackService := service.NewTrackService(trackRepo)

	consumer := mb.NewConsumer(cfg.KafkaBroker, mb.KafkaTopic, cfg.KafkaGroupID)
	kafkaCons := kafkaConsumer.NewKafkaConsumer(consumer)

	episodeHandler := handler.NewEpisodeHandler(episodeService)
	trackHandler := handler.NewTrackHandler(trackService)

	r := router.SetupRouter(episodeHandler, trackHandler, cfg.InternalAPIKey)
	r.SetTrustedProxies(nil)

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go kafkaCons.ConsumeEpisodes(ctx, episodeService.HandleEpisodeCreatedEvent)

	go func() {
		log.Println("archive-service started on :8080")
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
