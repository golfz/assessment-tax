package admin

import (
	"errors"
)

var (
	ErrReadingRequestBody       = errors.New("cannot reading request body")
	ErrInvalidInput             = errors.New("invalid input")
	ErrSettingPersonalDeduction = errors.New("error setting personal deduction")
)
