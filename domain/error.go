package domain

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrAlreadyTaken = errors.New("lock already taken")
	ErrUserExists   = errors.New("用户名已存在")
)

// 业务参数错误
type paramsError struct {
	err error
}

func (e *paramsError) Error() string {
	return e.err.Error()
}

func (e *paramsError) Unwrap() error {
	return e.err
}

func ParamsError(err error) error {
	return &paramsError{err: err}
}

func IsParamsError(err error) bool {
	if err == nil {
		return false
	}
	var e *paramsError
	return errors.As(err, &e)
}

func ParamsErrorf(format string, args ...any) error {
	return ParamsError(fmt.Errorf(format, args...))
}

// notFoundError 资源不存在错误
type notFoundError struct {
	err error
}

func (e *notFoundError) Unwrap() error {
	return e.err
}

func (e *notFoundError) Error() string {
	return e.err.Error()
}

func NotFoundError(err error) error {
	return &notFoundError{err: err}
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	var e *notFoundError
	return errors.As(err, &e)
}

// 冲突错误类型（用于唯一性约束等场景）
type conflictError struct {
	err error
}

func (e *conflictError) Error() string {
	return e.err.Error()
}

func (e *conflictError) Unwrap() error {
	return e.err
}

func ConflictError(err error) error {
	return &conflictError{err: err}
}

func IsConflict(err error) bool {
	if err == nil {
		return false
	}
	var e *conflictError
	return errors.As(err, &e)
}

// alreadyTakenError 锁已被占用错误
type alreadyTakenError struct {
	err error
}

func (e *alreadyTakenError) Error() string {
	return e.err.Error()
}

func (e *alreadyTakenError) Unwrap() error {
	return e.err
}

func IsAlreadyTakenError(err error) bool {
	if err == nil {
		return false
	}
	var e *alreadyTakenError
	return errors.As(err, &e)
}

func AlreadyTakenError(err error) error {
	return &alreadyTakenError{err: err}
}
