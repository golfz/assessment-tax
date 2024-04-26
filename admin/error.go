package admin

import (
	"errors"
)

var (
	ErrReadingRequestBody       = errors.New("cannot reading request body")
	ErrInvalidInput             = errors.New("invalid input")
	ErrInvalidPersonalDeduction = errors.New("invalid personal deduction")
	ErrSettingPersonalDeduction = errors.New("error setting personal deduction")
)
