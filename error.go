package goerror

import (
	"fmt"
	"net/http"
)

type GoError struct {
	Status int
	Code   string
	Msg    string
	Cause  string

	frames []*frame
}

func (e *GoError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e *GoError) IsCodeEqual(err error) bool {
	if ge, ok := err.(*GoError); ok {
		return ge.Code == e.Code
	}

	return false
}

func (e *GoError) WithCause(cause error) *GoError {
	e.Cause = cause.Error()
	e.frames = trace(DefaultStackTraceSkipLine)

	return e
}

func DefineBadRequest(code, msg string) *GoError {
	return &GoError{
		Status: http.StatusBadRequest,
		Code:   code,
		Msg:    msg,
	}
}

func DefineUnauthorized(code, msg string) *GoError {
	return &GoError{
		Status: http.StatusUnauthorized,
		Code:   code,
		Msg:    msg,
	}
}

func DefineForbidden(code, msg string) *GoError {
	return &GoError{
		Status: http.StatusForbidden,
		Code:   code,
		Msg:    msg,
	}
}

func DefineNotFound(code, msg string) *GoError {
	return &GoError{
		Status: http.StatusNotFound,
		Code:   code,
		Msg:    msg,
	}
}

func DefineInternalServerError(code, msg string) *GoError {
	return &GoError{
		Status: http.StatusInternalServerError,
		Code:   code,
		Msg:    msg,
	}
}

func DefineNotImplemented(code, msg string) *GoError {
	return &GoError{
		Status: http.StatusNotImplemented,
		Code:   code,
		Msg:    msg,
	}
}

func DefineBadGateway(code, msg string) *GoError {
	return &GoError{
		Status: http.StatusBadGateway,
		Code:   code,
		Msg:    msg,
	}
}

func DefineServiceUnavailable(code, msg string) *GoError {
	return &GoError{
		Status: http.StatusServiceUnavailable,
		Code:   code,
		Msg:    msg,
	}
}

func DefineGatewayTimeout(code, msg string) *GoError {
	return &GoError{
		Status: http.StatusGatewayTimeout,
		Code:   code,
		Msg:    msg,
	}
}
