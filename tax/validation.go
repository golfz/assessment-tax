package tax

import (
	"errors"
	"github.com/golfz/assessment-tax/rule"
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
	if deduction.Personal <= rule.MinPersonalDeduction || deduction.Personal > rule.MaxPersonalDeduction {
		err = errors.Join(err, ErrInvalidPersonalDeduction)
	}

	if deduction.KReceipt <= rule.MinKReceiptDeduction || deduction.KReceipt > rule.MaxKReceiptDeduction {
		err = errors.Join(err, ErrInvalidKReceiptDeduction)
	}

	if deduction.Donation > rule.MaxDonationDeduction {
		err = errors.Join(err, ErrInvalidDonationDeduction)
	}
	return
}
