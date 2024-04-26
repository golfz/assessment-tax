package tax

import (
	"github.com/golfz/assessment-tax/deduction"
	"strconv"
)

func CalculateTaxFromCSV(records [][]string, deductionData deduction.Deduction) (CsvTaxResponse, error) {
	result := CsvTaxResponse{}

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
					return CsvTaxResponse{}, ErrParsingData
				}
				taxInfo.TotalIncome = totalIncome
			case 1:
				wht, err := strconv.ParseFloat(colVal, 64)
				if err != nil {
					return CsvTaxResponse{}, ErrParsingData
				}
				taxInfo.WHT = wht
			case 2:
				donation, err := strconv.ParseFloat(colVal, 64)
				if err != nil {
					return CsvTaxResponse{}, ErrParsingData
				}
				taxInfo.Allowances = []Allowance{
					{Type: AllowanceTypeDonation, Amount: donation},
				}
			}
		}

		taxResult, err := CalculateTax(taxInfo, deductionData)
		if err != nil {
			return CsvTaxResponse{}, ErrCalculatingTax
		}

		result.Taxes = append(result.Taxes, CsvTaxRecord{
			TotalIncome: taxInfo.TotalIncome,
			Tax:         taxResult.Tax,
			TaxRefund:   taxResult.TaxRefund,
		})
	}

	return result, nil
}
