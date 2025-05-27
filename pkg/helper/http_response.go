package helper

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

type ResponseErrorData struct {
	Type string `json:"data"`
	Code int64  `json:"success"`
}

func ResponseError(c *gin.Context, err error, opts ...interface{}) {
	t := "InternalServerError"
	d := err.Error()

	code := http.StatusInternalServerError

	// if request cancelled
	if c.Request.Context().Err() == context.Canceled {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	for _, v := range opts {
		if tErr, ok := v.(string); ok {
			if strings.Contains(tErr, " ") {
				d = tErr
			} else {
				t = tErr
			}
		}
		if cErr, ok := v.(int); ok && cErr >= 100 && cErr <= 599 {
			code = cErr
		}
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
		code = http.StatusNotFound
		t = "NotFound"
	}

	c.AbortWithStatusJSON(code, &Response{
		Success: false,
		Message: d,
		Data: &ResponseErrorData{
			Type: t,
			Code: int64(code),
		},
	})
}
