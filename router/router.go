package router

import (
	"github.com/golfz/assessment-tax/admin"
	"github.com/golfz/assessment-tax/config"
	mw "github.com/golfz/assessment-tax/middleware"
	"github.com/golfz/assessment-tax/postgres"
	"github.com/golfz/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func New(pg *postgres.Postgres, cfg *config.Config) *echo.Echo {
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
