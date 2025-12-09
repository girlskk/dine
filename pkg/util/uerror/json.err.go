package uerror

import (
	"fmt"
	"strings"
)

type errUnmarshalJSON struct {
	raw []byte
	err error
	msg string
}

func (e *errUnmarshalJSON) Error() string {
	msg := e.msg
	if msg == "" {
		msg = "failed to unmarshaling JSON"
	}

	return fmt.Sprintf("%s: %q: %v", msg, e.raw, e.err)
}

func (e *errUnmarshalJSON) Unwrap() error {
	return e.err
}

func UnmarshalJSONErr(raw []byte, err error, msg ...string) error {
	return &errUnmarshalJSON{
		raw: raw,
		err: err,
		msg: strings.Join(msg, ":"),
	}
}

type errMarshalJSON struct {
	obj interface{}
	err error
	msg string
}

func (e *errMarshalJSON) Error() string {
	msg := e.msg
	if msg == "" {
		msg = "failed to marshaling JSON"
	}

	return fmt.Sprintf("%s: %+v: %v", msg, e.obj, e.err)
}

func (e *errMarshalJSON) Unwrap() error {
	return e.err
}

func MarshalJSONErr(obj interface{}, err error, msg ...string) error {
	return &errMarshalJSON{
		obj: obj,
		err: err,
		msg: strings.Join(msg, ":"),
	}
}
