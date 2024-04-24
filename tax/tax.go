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
	Amount float64       `json:"amount" validate:"min=0"`
}

type TaxInformation struct {
	TotalIncome float64     `json:"totalIncome" validate:"required,min=0"`
	WHT         float64     `json:"wht" validate:"min=0"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResult struct {
	Tax       float64 `json:"tax,omitempty"`
	TaxRefund float64 `json:"taxRefund,omitempty"`
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

	ErrInvalidPersonalDeduction = errors.New("invalid personal deduction")
	ErrInvalidKReceiptDeduction = errors.New("invalid k-receipt deduction")
	ErrInvalidDonationDeduction = errors.New("invalid donation deduction")

	ErrInvalidTaxInformation = errors.New("invalid tax information")

	ErrInvalidTotalIncome     = errors.New("total income must be greater than or equal to 0")
	ErrInvalidWHT             = errors.New("WHT must be greater than or equal to 0 and less than total income")
	ErrInvalidAllowanceAmount = errors.New("allowance amount must be greater than or equal to 0")
)

func validateTaxInformation(info TaxInformation) (err error) {
	if info.TotalIncome < 0 {
		err = errors.Join(err, ErrInvalidTotalIncome)
	}

	if info.WHT < 0 {
		err = errors.Join(err, ErrInvalidWHT)
	}

	if info.TotalIncome > 0 && info.WHT > info.TotalIncome {
		err = errors.Join(err, ErrInvalidWHT)
	}

	for _, allowance := range info.Allowances {
		if allowance.Amount < 0 {
			err = errors.Join(err, ErrInvalidAllowanceAmount)
		}
	}

	return
}

func validateDeduction(deduction Deduction) (err error) {
	if deduction.Personal <= ConstraintMinPersonalDeduction || deduction.Personal > ConstraintMaxPersonalDeduction {
		err = errors.Join(err, ErrInvalidPersonalDeduction)
	}

	if deduction.KReceipt <= ConstraintMinKReceiptDeduction || deduction.KReceipt > ConstraintMaxKReceiptDeduction {
		err = errors.Join(err, ErrInvalidKReceiptDeduction)
	}

	if deduction.Donation > ConstraintMaxDonationDeduction {
		err = errors.Join(err, ErrInvalidDonationDeduction)
	}
	return
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

	netIncome := info.TotalIncome - deduction.Personal

	// Calculate tax
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
