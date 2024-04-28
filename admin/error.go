package admin

import (
	"errors"
)

var (
	ErrReadingRequestBody       = errors.New("cannot reading request body")
	ErrInvalidInput             = errors.New("invalid input")
	ErrInvalidInputDeduction    = errors.New("invalid input deduction")
	ErrSettingPersonalDeduction = errors.New("error setting personal deduction")
	ErrSettingKReceiptDeduction = errors.New("error setting k-receipt deduction")
)
