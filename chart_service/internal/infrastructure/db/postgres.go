package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(databaseURL string) *pgxpool.Pool {

	for i := 0; i < 10; i++ {

		pool, err := pgxpool.New(context.Background(), databaseURL)
		if err == nil {

			err = pool.Ping(context.Background())
			if err == nil {
				log.Println("connected to postgres")
				return pool
			}
		}

		log.Println("waiting for postgres...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("cannot connect to postgres")
	return nil
}
