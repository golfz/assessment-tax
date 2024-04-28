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

type ProcessingInput struct {
	validateDeduction    ValidatorFunc
	setDeduction         SetterFunc
	errValidationInvalid error
	errSettingDeduction  error
}

func (h *Handler) DeductionProcessing(c echo.Context, processingInput ProcessingInput) (Deduction, int, error) {
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
	err = processingInput.validateDeduction(input.Deduction)
	if err != nil {
		c.Logger().Printf("error validating deduction: %v", err)
		return Deduction{}, http.StatusBadRequest, processingInput.errValidationInvalid
	}
	err = processingInput.setDeduction(input.Deduction)
	if err != nil {
		c.Logger().Printf("error setting deduction: %v", err)
		return Deduction{}, http.StatusInternalServerError, processingInput.errSettingDeduction
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
	data, statusCode, err := h.DeductionProcessing(c, ProcessingInput{
		validateDeduction:    deduction.ValidatePersonalDeduction,
		setDeduction:         h.store.SetPersonalDeduction,
		errValidationInvalid: ErrInvalidPersonalDeduction,
		errSettingDeduction:  ErrSettingPersonalDeduction,
	})
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

	err = deduction.ValidateKReceiptDeduction(input.Deduction)
	if err != nil {
		c.Logger().Printf("error validating k-receipt deduction: %v", err)
		return c.JSON(http.StatusBadRequest, Err{Message: ErrInvalidKReceiptDeduction.Error()})
	}

	err = h.store.SetKReceiptDeduction(input.Deduction)
	if err != nil {
		c.Logger().Printf("error setting k-receipt deduction: %v", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: ErrSettingKReceiptDeduction.Error()})
	}

	return c.JSON(http.StatusOK, KReceiptDeduction(input))
}
