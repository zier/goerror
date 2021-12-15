package goerror

import (
	"encoding/json"
	"fmt"
)

type Error interface {
	Error() string
	ErrorWithCause() string

	PrintRawJSONInput() string
	StackTrace() string
	Cause() string

	IsCodeEqual(err error) bool

	GetReasons() []*Reason
	AddReason(fieldName, reason string, value interface{})

	Input() map[string]interface{}

	WithCause(cause error) Error
	WithInput(input interface{}) Error
	WithKeyValueInput(inputs map[string]interface{}) Error
	WithExtendMsg(msg string) Error
}

type GoError struct {
	Status    int
	Code      string
	Msg       string
	ExtendMsg string
	cause     string

	reasons []*Reason

	input  map[string]interface{}
	frames []*frame
}

func (e *GoError) Error() string {
	if e.cause != "" {
		return fmt.Sprintf("%s: %s - %s", e.Code, e.Msg, e.cause)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e *GoError) PrintRawJSONInput() string {
	if e.input == nil {
		return ""
	}

	inputData, _ := json.Marshal(e.input)
	return string(inputData)
}

func (e *GoError) Input() map[string]interface{} {
	return e.input
}

func (e *GoError) Cause() string {
	return e.cause
}

func (e *GoError) ErrorWithCause() string {
	return fmt.Sprintf("%s - %s", e.Error(), e.Cause())
}

func (e *GoError) IsCodeEqual(err error) bool {
	if ge, ok := err.(*GoError); ok {
		return ge.Code == e.Code
	}

	return false
}

func (e *GoError) WithCause(cause error) Error {
	e.cause = cause.Error()
	e.frames = trace(DefaultStackTraceSkipLine)

	return e
}

func (e *GoError) WithInput(input interface{}) Error {
	if input == nil {
		e.input = map[string]interface{}{}
	} else {
		e.input = map[string]interface{}{
			"input": input,
		}
	}

	return e
}

func (e *GoError) WithExtendMsg(extendMsg string) Error {
	e.ExtendMsg = extendMsg

	return e
}

func (e *GoError) WithKeyValueInput(inputs map[string]interface{}) Error {
	if nil == inputs {
		e.input = map[string]interface{}{}
	} else {
		e.input = inputs
	}

	return e
}
