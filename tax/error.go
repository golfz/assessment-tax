package tax

import "errors"

var (
	ErrReadingRequestBody = errors.New("cannot reading request body")
	ErrGettingDeduction   = errors.New("error getting deduction")
	ErrCalculatingTax     = errors.New("error calculating tax")
)

// Invalid tax information errors
var (
	ErrInvalidTaxInformation = errors.New("invalid tax information")

	ErrInvalidTotalIncome     = errors.New("total income must be greater than or equal to 0")
	ErrInvalidWHT             = errors.New("WHT must be greater than or equal to 0 and less than total income")
	ErrInvalidAllowanceAmount = errors.New("allowance amount must be greater than or equal to 0")
)

// Invalid deduction errors
var (
	ErrInvalidDeduction = errors.New("invalid deduction")

	ErrInvalidPersonalDeduction = errors.New("invalid personal deduction")
	ErrInvalidKReceiptDeduction = errors.New("invalid k-receipt deduction")
	ErrInvalidDonationDeduction = errors.New("invalid donation deduction")
)
