package goerror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError_IsTypeEqual(t *testing.T) {
	err := DefineNotFound("NotFoundUser", "User not found")
	errUserNotFound := DefineNotFound("NotFoundUser", "User2 not found")

	require.True(t, err.IsCodeEqual(errUserNotFound))
	require.True(t, err.IsCodeEqual(err))
}

func TestError_IsTypeEqual_WithDefaultError(t *testing.T) {
	err := DefineBadRequest("InvalidRequest", "Username is required")
	errUnableGetStaff := errors.New("Unable get staff")

	require.False(t, err.IsCodeEqual(errUnableGetStaff))
}
