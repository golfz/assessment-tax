package tax

import (
	"errors"
	"math"
)

var rates = []rate{
	{lowerBound: 0, upperBound: 150_000, percentage: 0},
	{lowerBound: 150_000, upperBound: 500_000, percentage: 10},
	{lowerBound: 500_000, upperBound: 1_000_000, percentage: 15},
	{lowerBound: 1_000_000, upperBound: 2_000_000, percentage: 20},
	{lowerBound: 2_000_000, upperBound: math.MaxFloat64, percentage: 35},
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

func CalculateTax(info TaxInformation, deduction Deduction) (TaxResult, error) {
	err := validateTaxInformation(info)
	if err != nil {
		err = errors.Join(err, ErrInvalidTaxInformation)
		return TaxResult{}, err
	}

	err = validateDeduction(deduction)
	if err != nil {
		err = errors.Join(err, ErrInvalidDeduction)
		return TaxResult{}, err
	}

	totalAllowance := getTotalAllowance(info.Allowances, deduction)

	netIncome := calculateNetIncome(info.TotalIncome, deduction.Personal, totalAllowance)

	tax := 0.0
	for _, r := range rates {
		tax += calculateTaxForRate(r, netIncome)
	}

	tax -= info.WHT
	taxRefund := 0.0
	if tax < 0 {
		taxRefund = -tax
		tax = 0
	}

	return TaxResult{
		Tax:       tax,
		TaxRefund: taxRefund,
	}, nil
}
