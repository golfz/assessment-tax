package admin

import (
	"github.com/go-playground/validator/v10"
	"github.com/golfz/assessment-tax/deduction"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Storer interface {
	SetPersonalDeduction(amount float64) error
	SetKReceiptDeduction(amount float64) error
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

type ValidatorFunc func(float64) error
type SetterFunc func(float64) error

func (h *Handler) DeductionProcessing(c echo.Context, validateDeduction ValidatorFunc, setDeduction SetterFunc) (Deduction, int, error) {
	var input Deduction
	err := c.Bind(&input)
	if err != nil {
		c.Logger().Printf("error reading request body: %v", err)
		return Deduction{}, http.StatusBadRequest, ErrReadingRequestBody
	}
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		c.Logger().Printf("error validating request body: %v", err)
		return Deduction{}, http.StatusBadRequest, ErrInvalidInput
	}
	err = validateDeduction(input.Deduction)
	if err != nil {
		c.Logger().Printf("error validating deduction: %v", err)
		return Deduction{}, http.StatusBadRequest, ErrInvalidInputDeduction
	}
	err = setDeduction(input.Deduction)
	if err != nil {
		c.Logger().Printf("error setting deduction: %v", err)
		return Deduction{}, http.StatusInternalServerError, ErrSettingDeduction
	}
	return input, http.StatusOK, nil
}

// SetPersonalDeductionHandler
//
//	     @Security       BasicAuth
//			@Summary		Admin set personal deduction
//			@Description	Admin set personal deduction
//			@Tags			admin
//		    @Accept			json
//		    @Param			amount		body		Deduction	true		"Amount to set personal deduction"
//			@Produce		json
//			@Success		200	            {object}	PersonalDeduction
//			@Failure		400	            {object}	Err
//			@Failure		401	            {object}	Err
//			@Failure		500	            {object}	Err
//			@Router			/admin/deductions/personal [post]
func (h *Handler) SetPersonalDeductionHandler(c echo.Context) error {
	data, statusCode, err := h.DeductionProcessing(c, deduction.ValidatePersonalDeduction, h.store.SetPersonalDeduction)
	if err != nil {
		return c.JSON(statusCode, Err{Message: err.Error()})
	}
	return c.JSON(statusCode, PersonalDeduction(data))
}

// SetKReceiptDeductionHandler
//
//	         @Security       BasicAuth
//				@Summary		Admin set k-receipt deduction
//				@Description	Admin set k-receipt deduction
//				@Tags			admin
//			    @Accept			json
//			    @Param			amount		body		Deduction	true		"Amount to set personal deduction"
//				@Produce		json
//				@Success		200	            {object}	KReceiptDeduction
//				@Failure		400	            {object}	Err
//				@Failure		401	            {object}	Err
//				@Failure		500	            {object}	Err
//				@Router			/admin/deductions/k-receipt [post]
func (h *Handler) SetKReceiptDeductionHandler(c echo.Context) error {
	data, statusCode, err := h.DeductionProcessing(c, deduction.ValidateKReceiptDeduction, h.store.SetKReceiptDeduction)
	if err != nil {
		return c.JSON(statusCode, Err{Message: err.Error()})
	}
	return c.JSON(statusCode, KReceiptDeduction(data))
}
