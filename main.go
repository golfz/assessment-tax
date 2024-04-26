package main

import (
	"context"
	"crypto/subtle"
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

	_ "github.com/lib/pq"

	_ "github.com/golfz/assessment-tax/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title		K-Tax API
// @version		1.0
// @description Sophisticated K-Tax API
// @host		localhost:8080
// @BasePath    /
func main() {
	cfg := config.NewWith(os.Getenv)

	pg, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("exit: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	hTax := tax.New(pg)
	e.POST("/tax/calculations", hTax.CalculateTaxHandler)

	a := e.Group("/admin")

	a.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("adminTax")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("admin!")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	hAdmin := admin.New(pg)
	a.POST("/deductions/personal", hAdmin.SetPersonalDeductionHandler)

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
