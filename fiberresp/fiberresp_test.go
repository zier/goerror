package fiberresp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zier/goerror"
)

func TestRespWithError(t *testing.T) {
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		return RespWithError(c, goerror.DefineBadRequest("InvalidRequest", "Username is required"))
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodPost, "/", nil))
	assert.NoError(t, err)

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)

	assert.Equal(t, `{"errors":[],"message":"Username is required","type":"InvalidRequest"}`, string(data))
}
