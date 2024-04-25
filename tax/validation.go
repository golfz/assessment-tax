package tax

import "errors"

const (
	ConstraintMinPersonalDeduction float64 = 10_000.0
	ConstraintMaxPersonalDeduction float64 = 100_000.0

	ConstraintMaxDonationDeduction float64 = 100_000.0

	ConstraintMinKReceiptDeduction float64 = 0.0
	ConstraintMaxKReceiptDeduction float64 = 100_000.0
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
