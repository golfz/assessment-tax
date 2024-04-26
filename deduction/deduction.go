package deduction

import "errors"

const (
	MinPersonalDeduction float64 = 10_000.0
	MaxPersonalDeduction float64 = 100_000.0

	MaxDonationDeduction float64 = 100_000.0

	MinKReceiptDeduction float64 = 0.0
	MaxKReceiptDeduction float64 = 100_000.0
)

type Deduction struct {
	Personal float64
	KReceipt float64
	Donation float64
}

var (
	ErrInvalidPersonalDeduction = errors.New("invalid personal deduction")
	ErrInvalidKReceiptDeduction = errors.New("invalid k-receipt deduction")
	ErrInvalidDonationDeduction = errors.New("invalid donation deduction")
)

func ValidatePersonalDeduction(personal float64) (err error) {
	if personal <= MinPersonalDeduction || personal > MaxPersonalDeduction {
		err = errors.Join(err, ErrInvalidPersonalDeduction)
	}
	return
}

func ValidateKReceiptDeduction(kReceipt float64) (err error) {
	if kReceipt <= MinKReceiptDeduction || kReceipt > MaxKReceiptDeduction {
		err = errors.Join(err, ErrInvalidKReceiptDeduction)
	}
	return
}

func ValidateDonationDeduction(donation float64) (err error) {
	if donation > MaxDonationDeduction {
		err = errors.Join(err, ErrInvalidDonationDeduction)
	}
	return
}

func (d Deduction) Validate() (err error) {
	if e := ValidatePersonalDeduction(d.Personal); e != nil {
		err = errors.Join(err, e)
	}

	if e := ValidateKReceiptDeduction(d.KReceipt); e != nil {
		err = errors.Join(err, e)
	}

	if e := ValidateDonationDeduction(d.Donation); e != nil {
		err = errors.Join(err, e)
	}
	return
}
