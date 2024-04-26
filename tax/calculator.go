package tax

import (
	"errors"
	"github.com/golfz/assessment-tax/deduction"
	"math"
	"strconv"
)

var rates = []rate{
	{lowerBound: 0, upperBound: 150_000, percentage: 0, description: "0-150,000"},
	{lowerBound: 150_000, upperBound: 500_000, percentage: 10, description: "150,001-500,000"},
	{lowerBound: 500_000, upperBound: 1_000_000, percentage: 15, description: "500,001-1,000,000"},
	{lowerBound: 1_000_000, upperBound: 2_000_000, percentage: 20, description: "1,000,001-2,000,000"},
	{lowerBound: 2_000_000, upperBound: math.MaxFloat64, percentage: 35, description: "2,000,001 ขึ้นไป"},
}

func calculateTaxableIncome(netIncome, lowerBound, upperBound float64) float64 {
	if netIncome <= lowerBound {
		return 0
	}
	taxableIncome := netIncome - lowerBound
	if netIncome > upperBound {
		taxableIncome = upperBound - lowerBound
	}
	return taxableIncome
}

func calculateTaxForRate(r rate, netIncome float64) float64 {
	taxableIncome := calculateTaxableIncome(netIncome, r.lowerBound, r.upperBound)
	return taxableIncome * (r.percentage / 100.0)
}

func calculateNetIncome(totalIncome, personalDeduction, totalAllowance float64) float64 {
	result := totalIncome - personalDeduction - totalAllowance
	if result < 0 {
		return 0
	}
	return result
}

func CalculateTax(info TaxInformation, deduction deduction.Deduction) (TaxResult, error) {
	err := validateTaxInformation(info)
	if err != nil {
		err = errors.Join(err, ErrInvalidTaxInformation)
		return TaxResult{}, err
	}

	err = deduction.Validate()
	if err != nil {
		err = errors.Join(err, ErrInvalidDeduction)
		return TaxResult{}, err
	}

	totalAllowance := getTotalAllowance(info.Allowances, deduction)

	netIncome := calculateNetIncome(info.TotalIncome, deduction.Personal, totalAllowance)

	taxResult := TaxResult{
		Tax:       0.0,
		TaxRefund: 0.0,
		TaxLevels: make([]TaxLevel, 0),
	}
	for _, r := range rates {
		tax := calculateTaxForRate(r, netIncome)
		taxResult.Tax += tax
		taxResult.TaxLevels = append(taxResult.TaxLevels, TaxLevel{
			Level: r.description,
			Tax:   tax,
		})
	}

	taxResult.Tax -= info.WHT
	if taxResult.Tax < 0 {
		taxResult.TaxRefund = -taxResult.Tax
		taxResult.Tax = 0.0
	}

	return taxResult, nil
}

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
