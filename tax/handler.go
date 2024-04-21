package tax

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Storer interface {
}

type Handler struct {
	store Storer
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

func (h *Handler) CalculateTaxHandler(c echo.Context) error {
	t := TaxResult{Tax: 29000.0}
	return c.JSON(http.StatusOK, t)
}
