package validator

import (
	"github.com/gin-gonic/gin"
	gogrpcgatewayginresponse "github.com/ralvarezdev/go-grpc-gateway/gin/response"
	"google.golang.org/grpc/status"
	"net/http"
)

type (
	// DefaultHandler struct
	DefaultHandler struct {
		responseHandler gogrpcgatewayginresponse.Handler
	}
)

// NewDefaultHandler creates a new default response handler
func NewDefaultHandler(
	responseHandler gogrpcgatewayginresponse.Handler,
) (*DefaultHandler, error) {
	// Check if the response handler is nil
	if responseHandler == nil {
		return nil, gogrpcgatewayginresponse.ErrNilHandler
	}

	return &DefaultHandler{
		responseHandler: responseHandler,
	}, nil
}

// HandleError handles the error
func (d *DefaultHandler) HandleError(
	ctx *gin.Context,
	err error,
) {
	// Check if the error is a gRPC status error
	if _, ok := status.FromError(err); ok {
		d.responseHandler.HandleErrorResponse(ctx, err)
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	ctx.Abort()
}
