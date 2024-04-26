package admin

import (
	"github.com/go-playground/validator/v10"
	"github.com/golfz/assessment-tax/deduction"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Storer interface {
	SetPersonalDeduction(amount float64) error
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

// SetPersonalDeductionHandler
//
//		@Summary		Admin set personal deduction
//		@Description	Admin set personal deduction
//		@Tags			admin
//	    @Accept			json
//	    @Param			amount		body		Input	true		"Amount to set personal deduction"
//		@Produce		json
//		@Success		200	            {object}	PersonalDeduction
//		@Failure		400	            {object}	Err
//		@Failure		500	            {object}	Err
//		@Router			/admin/deductions/personal [post]
func (h *Handler) SetPersonalDeductionHandler(c echo.Context) error {
	var input Input
	err := c.Bind(&input)
	if err != nil {
		c.Logger().Printf("error reading request body: %v", err)
		return c.JSON(http.StatusBadRequest, Err{Message: ErrReadingRequestBody.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		c.Logger().Printf("error validating request body: %v", err)
		return c.JSON(http.StatusBadRequest, Err{Message: ErrInvalidInput.Error()})
	}

	err = deduction.ValidatePersonalDeduction(input.Amount)
	if err != nil {
		c.Logger().Printf("error validating personal deduction: %v", err)
		return c.JSON(http.StatusBadRequest, Err{Message: ErrInvalidPersonalDeduction.Error()})
	}

	err = h.store.SetPersonalDeduction(input.Amount)
	if err != nil {
		c.Logger().Printf("error setting personal deduction: %v", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: ErrSettingPersonalDeduction.Error()})
	}

	return c.NoContent(http.StatusOK)
}
