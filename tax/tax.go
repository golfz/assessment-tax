package tax

import (
	"errors"
	"math"
)

type AllowanceType string

const (
	AllowanceTypeDonation AllowanceType = "donation"
	AllowanceTypeKReceipt AllowanceType = "k-receipt"

	ConstraintMinPersonalDeduction float64 = 10_000.0
	ConstraintMaxPersonalDeduction float64 = 100_000.0

	ConstraintMaxDonationDeduction float64 = 100_000.0

	ConstraintMinKReceiptDeduction float64 = 0.0
	ConstraintMaxKReceiptDeduction float64 = 100_000.0
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
	KReceipt float64
	Donation float64
}

type rate struct {
	lowerBound float64
	upperBound float64
	percentage float64
}

var rates = []rate{
	{lowerBound: 0, upperBound: 150_000, percentage: 0},
	{lowerBound: 150_000, upperBound: 500_000, percentage: 10},
	{lowerBound: 500_000, upperBound: 1_000_000, percentage: 15},
	{lowerBound: 1_000_000, upperBound: 2_000_000, percentage: 20},
	{lowerBound: 2_000_000, upperBound: math.MaxFloat64, percentage: 35},
}

var (
	ErrInvalidDeduction = errors.New("invalid deduction")
)

func validateDeduction(deduction Deduction) error {
	if deduction.Personal <= ConstraintMinPersonalDeduction || deduction.Personal > ConstraintMaxPersonalDeduction {
		return ErrInvalidDeduction
	}

	if deduction.KReceipt <= ConstraintMinKReceiptDeduction || deduction.KReceipt > ConstraintMaxKReceiptDeduction {
		return ErrInvalidDeduction
	}

	if deduction.Donation > ConstraintMaxDonationDeduction {
		return ErrInvalidDeduction
	}
	return nil
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

func CalculateTax(info TaxInformation, deduction Deduction) (TaxResult, error) {
	err := validateDeduction(deduction)
	if err != nil {
		return TaxResult{}, err
	}

	netIncome := info.TotalIncome - deduction.Personal

	// Calculate tax
	tax := 0.0

	for _, r := range rates {
		tax += calculateTaxForRate(r, netIncome)
	}

	return TaxResult{
		Tax: tax,
	}, nil
}
