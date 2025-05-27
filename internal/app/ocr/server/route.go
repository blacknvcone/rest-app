package ocr

import (
	"rest-app/internal/app/ocr/port"

	"github.com/gin-gonic/gin"
)

type (
	routes struct{}
)

var (
	Routes routes
)

func (r routes) New(router *gin.RouterGroup, handler port.IOCRHandler) {
	router.POST("/receipt", handler.ProcessReceipt)
}
