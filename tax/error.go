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

// Deduction errors
var (
	ErrInvalidDeduction = errors.New("invalid deduction")
)

var (
	ErrUploadingFile    = errors.New("cannot uploading file")
	ErrReadingFile      = errors.New("cannot reading file")
	ErrReadingCSV       = errors.New("cannot reading csv")
	ErrParsingData      = errors.New("cannot parsing data")
	ErrInvalidCSVHeader = errors.New("invalid csv header")
)
