package goerror

import (
    "errors"
    "fmt"
    "io/ioutil"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestGoError_IsCodeEqual(t *testing.T) {
    err := DefineNotFound("NotFoundUser", "User not found")
    errUserNotFound := DefineNotFound("NotFoundUser", "User2 not found")

    require.True(t, err.IsCodeEqual(errUserNotFound))
    require.True(t, err.IsCodeEqual(err))
}

func TestGoError_IsCodeEqual_WithDefaultError(t *testing.T) {
    err := DefineBadRequest("InvalidRequest", "Username is required")
    errUnableGetStaff := errors.New("Unable get staff")

    require.False(t, err.IsCodeEqual(errUnableGetStaff))
}

func TestGoError_WithCause(t *testing.T) {
    _, err := ioutil.ReadFile("/tmp/dat")

    goErr := DefineInternalServerError("TestStackTrace", "Test stacktrace").WithCause(err)
    fmt.Println(goErr.StackTrace())
    require.Error(t, goErr)
}
