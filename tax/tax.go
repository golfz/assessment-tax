package tax

import (
	"math"
)

type AllowanceType string

const (
	AllowanceTypeDonation AllowanceType = "donation"
	AllowanceTypeKReceipt AllowanceType = "k-receipt"
)

type Allowance struct {
	Type   AllowanceType `json:"allowanceType"`
	Amount float64       `json:"amount"`
}

type TaxInformation struct {
	TotalIncome float64     `json:"totalIncome" validate:"required,min=0"`
	WHT         float64     `json:"wht" validate:"min=0"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResult struct {
	Tax float64 `json:"tax"`
}

type Deduction struct {
	Personal float64
}

type rate struct {
	moreThan   float64
	to         float64
	percentage float64
}

var rates = []rate{
	{moreThan: 0, to: 150_000, percentage: 0},
	{moreThan: 150_000, to: 500_000, percentage: 10},
	{moreThan: 500_000, to: 1_000_000, percentage: 15},
	{moreThan: 1_000_000, to: 2_000_000, percentage: 20},
	{moreThan: 2_000_000, to: math.MaxFloat64, percentage: 35},
}

func cal(r rate, netIncome float64) float64 {
	if netIncome > r.moreThan {
		taxableIncome := netIncome - r.moreThan
		taxRange := r.to - r.moreThan
		if taxableIncome > taxRange {
			taxableIncome = r.to - r.moreThan
		}
		return taxableIncome * (r.percentage / 100.0)
	}
	return 0
}

func CalculateTax(info TaxInformation, deduction Deduction) (TaxResult, error) {
	netIncome := info.TotalIncome - deduction.Personal

	// Calculate tax
	tax := 0.0

	for _, r := range rates {
		tax += cal(r, netIncome)
	}

	return TaxResult{
		Tax: tax,
	}, nil
}
