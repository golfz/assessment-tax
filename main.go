package main

import (
	"context"
	"fmt"
	"github.com/golfz/assessment-tax/config"
	"github.com/golfz/assessment-tax/postgres"
	"github.com/golfz/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	_ "github.com/golfz/assessment-tax/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	cfg := config.NewWith(os.Getenv)

	pg, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("exit: %v", err)
	}

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	hTax := tax.New(pg)
	e.POST("/tax/calculations", hTax.CalculateTaxHandler)

	// monitor shutdown signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	fmt.Println("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("server gracefully stopped")
}
