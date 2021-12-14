package fiberresp

import (
	"net/http"

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
