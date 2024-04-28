package admin

import (
	"errors"
)

var (
	ErrReadingRequestBody    = errors.New("cannot reading request body")
	ErrInputValidation       = errors.New("invalid input")
	ErrInvalidInputDeduction = errors.New("invalid input deduction")
	ErrSettingDeduction      = errors.New("error setting deduction")
)
