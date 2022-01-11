package fiberresp

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zier/goerror"
)

func RespWithError(c *fiber.Ctx, err error) error {
	if e, ok := err.(*goerror.GoError); ok {
		return c.Status(e.Status).JSON(fiber.Map{
			"type":    e.Code,
			"message": e.Msg + e.ExtendMsg,
			"errors":  e.GetReasons(),
		})
	}

	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"type":    "UnknownType",
		"message": err.Error(),
	})
}

func RespValidateError(c *fiber.Ctx, err error) error {
	errValidates := make([]*goerror.Reason, 0)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, errValidate := range errs {
			errValidates = append(errValidates, &goerror.Reason{
				FieldName: errValidate.Field(),
				Reason:    errValidate.ActualTag(),
				Value:     errValidate.Param(),
			})
		}
	}

	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"type":    "InvalidRequest",
		"message": "invalid request",
		"errors":  errValidates,
	})
}
