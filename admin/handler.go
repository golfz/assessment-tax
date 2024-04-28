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

func validateInput(c echo.Context, input *Deduction) (err error) {
	err = c.Bind(input)
	if err != nil {
		return ErrReadingRequestBody
	}
	validate := validator.New()
	if err = validate.Struct(*input); err != nil {
		return ErrInputValidation
	}
	return
}

func handleError(c echo.Context, errStatus int, err error, action, errMsg string) error {
	c.Logger().Printf("error %s: %v", action, err)
	return c.JSON(errStatus, Err{Message: errMsg})
}

func (h *Handler) DeductionProcessing(c echo.Context, validateDeduction ValidatorFunc, setDeduction SetterFunc, output OutputFunc) error {
	var input Deduction
	if err := validateInput(c, &input); err != nil {
		return handleError(c, http.StatusBadRequest, err, "reading request body", err.Error())
	}
	if err := validateDeduction(input.Deduction); err != nil {
		return handleError(c, http.StatusBadRequest, err, "validating deduction", ErrInvalidInputDeduction.Error())
	}
	if err := setDeduction(input.Deduction); err != nil {
		return handleError(c, http.StatusInternalServerError, err, "setting deduction", ErrSettingDeduction.Error())
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
//	     @Security       BasicAuth
//			@Summary		Admin set k-receipt deduction
//			@Description	Admin set k-receipt deduction
//			@Tags			admin
//			@Accept			json
//			@Param			amount		body		Deduction	true		"Amount to set personal deduction"
//			@Produce		json
//			@Success		200	            {object}	KReceiptDeduction
//			@Failure		400	            {object}	Err
//			@Failure		401	            {object}	Err
//			@Failure		500	            {object}	Err
//			@Router			/admin/deductions/k-receipt [post]
func (h *Handler) SetKReceiptDeductionHandler(c echo.Context) error {
	return h.DeductionProcessing(c, deduction.ValidateKReceiptDeduction, h.store.SetKReceiptDeduction, outputToKReceiptDeduction)
}
