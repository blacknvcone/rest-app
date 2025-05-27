package port

import "github.com/gin-gonic/gin"

type IOCRHandler interface {
	ProcessReceipt(ctx *gin.Context)
}
