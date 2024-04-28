package main

import (
	"context"
	"fmt"
	"github.com/golfz/assessment-tax/admin"
	"github.com/golfz/assessment-tax/config"
	"github.com/golfz/assessment-tax/postgres"
	"github.com/golfz/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mw "github.com/golfz/assessment-tax/middleware"

	_ "github.com/lib/pq"

	_ "github.com/golfz/assessment-tax/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func echoSetup(pg *postgres.Postgres, cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	hTax := tax.New(pg)
	e.POST("/tax/calculations", hTax.CalculateTaxHandler)
	e.POST("/tax/calculations/upload-csv", hTax.UploadCSVHandler)

	a := e.Group("/admin")
	a.Use(middleware.BasicAuth(mw.BasicAuth(*cfg)))

	hAdmin := admin.New(pg)
	a.POST("/deductions/personal", hAdmin.SetPersonalDeductionHandler)
	a.POST("/deductions/k-receipt", hAdmin.SetKReceiptDeductionHandler)

	return e
}

func startServer(e *echo.Echo, cfg *config.Config) {
	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
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

// @title		K-Tax API
// @version		1.0
// @description This is an API for K-Tax.
// @host		localhost:8080
// @BasePath    /
// @securityDefinitions.basic BasicAuth
func main() {
	cfg := config.NewWith(os.Getenv)

	pg, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("exit: %v", err)
	}

	e := echoSetup(pg, cfg)

	// monitor shutdown signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go startServer(e, cfg)

	waitForShutdown(ctx, e)
}
