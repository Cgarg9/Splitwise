package main

import (
	"context"
	"splitwise-clone/db"
	"splitwise-clone/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	logger.SetupGlobal("info")
	log.Info().Msg("Application started")

	ctx := context.Background()
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database pool")
	}
	defer pool.Close()

	log.Info().Msg("Database pool created successfully")
}
