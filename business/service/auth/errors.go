package auth

import (
	berrors2 "github.com/FlameInTheDark/rebot/foundation/berrors"
)

const baseCode = 10000

var (
	ErrInvalidInput = &berrors2.BusinessError{
		ErrCode: baseCode + 1,
		Message: "got invalid input",
	}
)
