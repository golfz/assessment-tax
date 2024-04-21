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

type Err struct {
	Message string `json:"message"`
}

// CalculateTaxHandler
//
//		@Summary		Calculate tax
//		@Description	Calculate tax
//		@Tags			tax
//	    @Accept			json
//	    @Param			amount		body		TaxInformation	true		"Amount to calculate tax"
//		@Produce		json
//		@Success		200	            {object}	TaxResult
//		@Failure		400	            {object}	Err
//		@Failure		500	            {object}	Err
//		@Router			/tax/calculations [post]
func (h *Handler) CalculateTaxHandler(c echo.Context) error {
	t := TaxResult{Tax: 29000.0}
	return c.JSON(http.StatusOK, t)
}
