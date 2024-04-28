package admin

import (
	"errors"
)

var (
	ErrReadingRequestBody    = errors.New("cannot reading request body")
	ErrInvalidInput          = errors.New("invalid input")
	ErrInvalidInputDeduction = errors.New("invalid input deduction")
	ErrSettingDeduction      = errors.New("error setting deduction")
)
