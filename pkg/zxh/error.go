package zxh

import "errors"

type APIError struct {
	ErrCode int
	ErrMsg  string
	ErrDlt  string
}

func (e *APIError) Error() string {
	return e.ErrMsg
}

func IsAPIError(err error) bool {
	if err == nil {
		return false
	}
	var e *APIError
	return errors.As(err, &e)
}
