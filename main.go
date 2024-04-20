package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golfz/assessment-tax/config"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("unable to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

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
