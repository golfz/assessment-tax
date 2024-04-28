package main

import (
	"context"
	"fmt"
	"github.com/golfz/assessment-tax/config"
	"github.com/golfz/assessment-tax/postgres"
	"github.com/golfz/assessment-tax/router"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golfz/assessment-tax/docs"
)

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

func initPostgres(cfg *config.Config) *postgres.Postgres {
	pg, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("exit: %v", err)
	}
	return pg
}

func startServer(e *echo.Echo, cfg *config.Config) {
	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}

func monitorShutdownSignal() (ctx context.Context, stop context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

func waitForShutdown(ctx context.Context, e *echo.Echo) {
	<-ctx.Done()
	fmt.Println("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("server gracefully stopped")
}
