package goerror

import (
    "fmt"
)

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

func (e *GoError) WithCause(cause error) *GoError {
    e.Cause = cause.Error()
    e.frames = trace(DefaultStackTraceSkipLine)

    return e
}

func (e *GoError) WithInput(input interface{}) *GoError {
    e.input = input

    return e
}
