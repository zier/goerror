package ginresp

import (
	"net/http"

	"github.com/devit-tel/goerror"
	"github.com/gin-gonic/gin"
)

func RespWithError(c *gin.Context, err error) {
	if e, ok := err.(*goerror.GoError); ok {
		c.JSON(e.Status, gin.H{
			"type":    e.Code,
			"message": e.Msg,
		})

		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"type":    "UnknownType",
		"message": err.Error(),
	})
}
