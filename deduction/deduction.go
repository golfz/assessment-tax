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

func (d Deduction) Validate() (err error) {
	if d.Personal <= MinPersonalDeduction || d.Personal > MaxPersonalDeduction {
		err = errors.Join(err, ErrInvalidPersonalDeduction)
	}

	if d.KReceipt <= MinKReceiptDeduction || d.KReceipt > MaxKReceiptDeduction {
		err = errors.Join(err, ErrInvalidKReceiptDeduction)
	}

	if d.Donation > MaxDonationDeduction {
		err = errors.Join(err, ErrInvalidDonationDeduction)
	}
	return
}
