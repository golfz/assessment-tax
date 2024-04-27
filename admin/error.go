package admin

import (
	"errors"
)

var (
	ErrReadingRequestBody       = errors.New("cannot reading request body")
	ErrInvalidInput             = errors.New("invalid input")
	ErrInvalidPersonalDeduction = errors.New("invalid personal deduction")
	ErrInvalidKReceiptDeduction = errors.New("invalid k-receipt deduction")
	ErrSettingPersonalDeduction = errors.New("error setting personal deduction")
	ErrSettingKReceiptDeduction = errors.New("error setting k-receipt deduction")
)
