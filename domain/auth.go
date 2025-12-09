package domain

import "errors"

var (
	ErrTokenInvalid = errors.New("token is invalid")
)

type AuthConfig struct {
	Expire int
	Secret string
}
