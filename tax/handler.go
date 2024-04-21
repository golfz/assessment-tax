package tax

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Storer interface {
	GetDeduction() (Deduction, error)
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
	var taxInfo TaxInformation
	err := c.Bind(&taxInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "bad request body"})
	}

	validate := validator.New()
	if err := validate.Struct(taxInfo); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "bad request body"})
	}

	deduction, err := h.store.GetDeduction()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "error getting deduction"})
	}

	result, err := CalculateTax(taxInfo, deduction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "error calculating tax"})
	}

	return c.JSON(http.StatusOK, result)
}
