package main

import (
	"log"

	"charts-user-service/internal/config"
	"charts-user-service/internal/infrastructure/cache"
	"charts-user-service/internal/infrastructure/db"
	"charts-user-service/internal/repository/postgres"
	"charts-user-service/internal/repository/redis"
	"charts-user-service/internal/service"
	"charts-user-service/internal/transport/http/router"
)

func main() {

	cfg := config.Load()

	pg := db.NewPostgresPool(cfg.DatabaseURL)
	db.RunMigrations(cfg.DatabaseURL)
	rdb, err := cache.NewRedisClient(cfg.RedisAddr)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepository(pg)
	authRepo := redis.NewAuthRepository(rdb)

	authService := service.NewAuthService(userRepo, authRepo)
	userService := service.NewUserService(userRepo)

	r := router.SetupRouter(authService, userService, cfg.InternalAPIKey)
	r.SetTrustedProxies(nil)
	log.Println(" user-service started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
