package tax

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golfz/assessment-tax/deduction"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Storer interface {
	GetDeduction() (deduction.Deduction, error)
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
		c.Logger().Printf("error reading request body: %v", err)
		return c.JSON(http.StatusBadRequest, Err{Message: ErrReadingRequestBody.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(taxInfo); err != nil {
		c.Logger().Printf("error validating request body: %v", err)
		return c.JSON(http.StatusBadRequest, Err{Message: ErrInvalidTaxInformation.Error()})
	}

	deductionData, err := h.store.GetDeduction()
	if err != nil {
		c.Logger().Printf("error getting deduction: %v", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: ErrGettingDeduction.Error()})
	}

	result, err := CalculateTax(taxInfo, deductionData)
	if err != nil {
		c.Logger().Printf("error calculating tax: %v", err)
		if errors.Is(err, ErrInvalidTaxInformation) {
			return c.JSON(http.StatusBadRequest, Err{Message: ErrInvalidTaxInformation.Error()})
		}
		return c.JSON(http.StatusInternalServerError, Err{Message: ErrCalculatingTax.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// UploadCSVHandler
//
//		@Summary		Upload csv file and calculate tax
//		@Description	Upload csv file and calculate tax
//		@Tags			tax
//	    @Accept			multipart/form-data
//	    @Param			taxFile	formData	file			true	"this is a test file"
//		@Produce		json
//		@Success		200	            {object}	CsvTaxResponse
//		@Failure		400	            {object}	Err
//		@Failure		500	            {object}	Err
//		@Router			/tax/calculations/upload-csv [post]
func (h *Handler) UploadCSVHandler(c echo.Context) error {
	file, err := c.FormFile("taxFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: ErrUploadingFile.Error()})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: ErrUploadingFile.Error()})
	}
	defer src.Close()

	cr := NewCSVReader(src)
	records, err := cr.readRecords()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: ErrReadingCSV.Error()})
	}

	deductionData, err := h.store.GetDeduction()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: ErrGettingDeduction.Error()})
	}

	result, err := CalculateTaxFromCSV(records, deductionData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: ErrCalculatingTax.Error()})

	}

	return c.JSON(http.StatusOK, result)
}
