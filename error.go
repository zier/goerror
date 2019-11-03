package goerror

import (
	"fmt"
)

type Error interface {
	Error() string
	PrintInput() string
	IsCodeEqual(err error) bool
	WithCause(cause error) Error
	WithInput(input interface{}) Error
	ExtendMsg(msg string) Error
}

type GoError struct {
	Status int
	Code   string
	Msg    string
	Cause  string

	input  interface{}
	frames []*frame
}

func (e *GoError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e *GoError) PrintInput() string {
	if e.input == nil {
		return ""
	}

	return fmt.Sprintf("%v", e.input)
}

func (e *GoError) IsCodeEqual(err error) bool {
	if ge, ok := err.(*GoError); ok {
		return ge.Code == e.Code
	}

	return false
}

func (e *GoError) WithCause(cause error) Error {
	e.Cause = cause.Error()
	e.frames = trace(DefaultStackTraceSkipLine)

	return e
}

func (e *GoError) WithInput(input interface{}) Error {
	e.input = input

	return e
}

func (e *GoError) ExtendMsg(msg string) Error {
	e.Msg = e.Msg + msg

	return e
}
