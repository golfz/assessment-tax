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
type OutputFunc func(Deduction) interface{}

func outputToPersonalDeduction(input Deduction) interface{} {
	return PersonalDeduction(input)
}

func outputToKReceiptDeduction(input Deduction) interface{} {
	return KReceiptDeduction(input)
}

func (h *Handler) DeductionProcessing(c echo.Context, validateDeduction ValidatorFunc, setDeduction SetterFunc, output OutputFunc) error {
	var input Deduction
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
	err = validateDeduction(input.Deduction)
	if err != nil {
		c.Logger().Printf("error validating deduction: %v", err)
		return c.JSON(http.StatusBadRequest, Err{Message: ErrInvalidInputDeduction.Error()})
	}
	err = setDeduction(input.Deduction)
	if err != nil {
		c.Logger().Printf("error setting deduction: %v", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: ErrSettingDeduction.Error()})
	}
	return c.JSON(http.StatusOK, output(input))
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
	return h.DeductionProcessing(c, deduction.ValidatePersonalDeduction, h.store.SetPersonalDeduction, outputToPersonalDeduction)
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
	return h.DeductionProcessing(c, deduction.ValidateKReceiptDeduction, h.store.SetKReceiptDeduction, outputToKReceiptDeduction)
}
