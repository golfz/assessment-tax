package tax

import (
	"errors"
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
