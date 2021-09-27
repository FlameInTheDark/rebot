package berrors

import (
	"strconv"

	"github.com/pkg/errors"
)

type BusinessError struct {
	ErrCode int
	Message string
}

func (e *BusinessError) Code() int {
	return e.ErrCode
}

func (e *BusinessError) Error() string {
	return e.Message + ": error code " + strconv.Itoa(e.ErrCode)
}

// WrapWithError wraps business error with another error
func WrapWithError(bErr *BusinessError, err error) error {
	return errors.Wrap(bErr, err.Error())
}
