package ginresp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devit-tel/goerror"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRespWithError_WithAftError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)
	res := httptest.NewRecorder()

	gin.SetMode("release")
	router := gin.New()
	router.GET("/user", func(c *gin.Context) {
		RespWithError(c, goerror.DefineBadRequest("InvalidRequest", "Username is required"))
	})

	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, `{"message":"Username is required","type":"InvalidRequest"}`, string(res.Body.Bytes()))
}

func TestRespWithError_WithDefaultError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)
	res := httptest.NewRecorder()

	gin.SetMode("release")
	router := gin.New()
	router.GET("/user", func(c *gin.Context) {
		RespWithError(c, errors.New("default error"))
	})

	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, `{"message":"default error","type":"UnknownType"}`, string(res.Body.Bytes()))
}
