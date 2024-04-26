package tax

import (
	"encoding/csv"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golfz/assessment-tax/deduction"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

	deduction, err := h.store.GetDeduction()
	if err != nil {
		c.Logger().Printf("error getting deduction: %v", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: ErrGettingDeduction.Error()})
	}

	result, err := CalculateTax(taxInfo, deduction)
	if err != nil {
		c.Logger().Printf("error calculating tax: %v", err)
		if errors.Is(err, ErrInvalidTaxInformation) {
			return c.JSON(http.StatusBadRequest, Err{Message: ErrInvalidTaxInformation.Error()})
		}
		return c.JSON(http.StatusInternalServerError, Err{Message: ErrCalculatingTax.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) UploadCSVHandler(c echo.Context) error {
	formFile, err := c.FormFile("taxFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: ErrUploadingFile.Error()})
	}
	file, err := formFile.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: ErrReadingFile.Error()})
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: ErrReadingCSV.Error()})
	}

	result := CSVResponse{}

	for i, row := range records {
		if i == 0 {
			continue
		}

		taxInfo := TaxInformation{}

		for colNum, colVal := range row {
			switch colNum {
			case 0:
				totalIncome, err := strconv.ParseFloat(colVal, 64)
				if err != nil {
					return c.JSON(http.StatusBadRequest, Err{Message: ErrParsingData.Error()})
				}
				taxInfo.TotalIncome = totalIncome
			case 1:
				wht, err := strconv.ParseFloat(colVal, 64)
				if err != nil {
					return c.JSON(http.StatusBadRequest, Err{Message: ErrParsingData.Error()})
				}
				taxInfo.WHT = wht
			case 2:
				donation, err := strconv.ParseFloat(colVal, 64)
				if err != nil {
					return c.JSON(http.StatusBadRequest, Err{Message: ErrParsingData.Error()})
				}
				taxInfo.Allowances = []Allowance{
					{Type: AllowanceTypeDonation, Amount: donation},
				}
			}
		}

		deductionData, err := h.store.GetDeduction()
		if err != nil {
			c.Logger().Printf("error getting deduction: %v", err)
			return c.JSON(http.StatusInternalServerError, Err{Message: ErrGettingDeduction.Error()})
		}

		taxResult, err := CalculateTax(taxInfo, deductionData)
		if err != nil {
			c.Logger().Printf("error calculating tax: %v", err)
			return c.JSON(http.StatusInternalServerError, Err{Message: ErrCalculatingTax.Error()})
		}

		result.Taxes = append(result.Taxes, CSVTaxResult{
			TotalIncome: taxInfo.TotalIncome,
			Tax:         taxResult.Tax,
			TaxRefund:   taxResult.TaxRefund,
		})

	}

	return c.JSON(http.StatusOK, result)
}
