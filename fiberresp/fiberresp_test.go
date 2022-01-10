package fiberresp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func newRequestWithBody(t *testing.T, jsonData interface{}) *http.Request {
	data, err := json.Marshal(jsonData)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	return req
}

func cleanResponse(responseBytes []byte) string {
	return strings.Replace(string(responseBytes), "\n", "", -1)
}

func TestRespValidateError(t *testing.T) {
	expectErrorJson := `{"errors":[{"fieldName":"age","reason":"lte","value":"130"}],"message":"invalid request","type":"InvalidRequest"}`

	reqData := struct {
		UserID string `json:"userID"`
		Name   string `json:"name"`
		Age    int    `json:"age"`
	}{
		UserID: "data1",
		Name:   "tester",
		Age:    500,
	}

	emptyStruct := struct {
		UserID string `json:"userID"`
		Name   string `json:"name" validate:"required"`
		Age    int    `json:"age" validate:"gte=0,lte=130"`
	}{}

	req := newRequestWithBody(t, reqData)

	app := fiber.New()
	app.Post("/user", func(c *fiber.Ctx) error {
		if err := c.BodyParser(&emptyStruct); err != nil {
			return err
		}

		validate := validator.New()
		// https://github.com/go-playground/validator/issues/258
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})

		errors := validate.Struct(emptyStruct)
		if errors != nil {
			RespValidateError(c, errors)
		}

		return nil
	})

	res, err := app.Test(req)
	require.NoError(t, err)
	// require.Equal(t, http.StatusBadRequest, res.StatusCode)

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	require.Equal(t, expectErrorJson, cleanResponse(data))
}

func TestRespWithErrorReasons(t *testing.T) {
	expectErrorJson := `{"errors":[{"fieldName":"username","reason":"username already exist"},{"fieldName":"phone","reason":"phone number already exist","value":"0598881111"}],"message":"user is already exist","type":"UserExist"}`

	req := newRequestWithBody(t, nil)

	app := fiber.New()
	app.Post("/user", func(c *fiber.Ctx) error {
		e := goerror.DefineBadRequest("UserExist", "user is already exist")
		e.AddReason("username", "username already exist", nil)
		e.AddReason("phone", "phone number already exist", "0598881111")

		return RespWithError(c, e)
	})

	res, err := app.Test(req)
	assert.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	require.Equal(t, expectErrorJson, cleanResponse(data))
}
