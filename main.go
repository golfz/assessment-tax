package main

import (
	"github.com/golfz/assessment-tax/config"
	"github.com/golfz/assessment-tax/postgres"
	"github.com/golfz/assessment-tax/router"
	_ "github.com/lib/pq"
	"log"
	"os"

	_ "github.com/golfz/assessment-tax/docs"
)

func initPostgres(cfg *config.Config) *postgres.Postgres {
	pg, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("exit: %v", err)
	}
	return pg
}

// @title		K-Tax API
// @version		1.0
// @description This is an API for K-Tax.
// @host		localhost:8080
// @BasePath    /
// @securityDefinitions.basic BasicAuth
func main() {
	cfg := config.NewWith(os.Getenv)
	pg := initPostgres(cfg)
	e := router.New(pg, cfg)

	ctx, stop := monitorShutdownSignal()
	defer stop()

	go startServer(e, cfg)
	waitForShutdown(ctx, e)
}
