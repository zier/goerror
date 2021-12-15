package goerror

import (
	"errors"
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
	require.Error(t, goErr)
	require.NotEmpty(t, goErr.StackTrace())
	require.Equal(t, "TestStackTrace: Test stacktrace - open /tmp/dat: no such file or directory - open /tmp/dat: no such file or directory", goErr.ErrorWithCause())
}

func TestGoError_Input(t *testing.T) {
	inputNil := DefineInternalServerError("TestInput", "Test input").WithInput(nil)
	require.Equal(t, "{}", inputNil.PrintRawJSONInput())

	inputString := DefineInternalServerError("TestInput", "Test input").WithInput("i am string")
	require.Equal(t, "{\"input\":\"i am string\"}", inputString.PrintRawJSONInput())

	inputStrings := DefineInternalServerError("TestInput", "Test input").WithInput([]string{"one", "two", "three"})
	require.Equal(t, "{\"input\":[\"one\",\"two\",\"three\"]}", inputStrings.PrintRawJSONInput())

	inputMap := DefineInternalServerError("TestInput", "Test input").WithInput(
		struct {
			UserID string `json:"userID"`
			Name   string `json:"name"`
		}{
			UserID: "user_1",
			Name:   "tester",
		})
	require.Equal(t, "{\"input\":{\"userID\":\"user_1\",\"name\":\"tester\"}}", inputMap.PrintRawJSONInput())
}

func TestGoError_WithKeyValueInput(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		err := DefineInternalServerError("TestInput", "Test input").WithKeyValueInput(
			map[string]interface{}{
				"job":   "knight",
				"level": 99,
				"desc":  "any",
			},
		)

		inputData := err.Input()

		require.Equal(t, inputData["job"], "knight")
		require.Equal(t, inputData["level"], 99)
		require.Equal(t, inputData["desc"], "any")
	})

	t.Run("Nil input", func(t *testing.T) {
		err := DefineInternalServerError("TestInput", "Test input").WithKeyValueInput(nil)
		require.Equal(t, map[string]interface{}{}, err.Input())

		inputData := err.Input()
		require.Equal(t, map[string]interface{}{}, inputData)
	})
}
