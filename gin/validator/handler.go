package validator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	goginresponse "github.com/ralvarezdev/go-gin/response"
	gogrpcgatewayginresponse "github.com/ralvarezdev/go-grpc-gateway/gin/response"
	"google.golang.org/grpc/status"
)

type (
	// DefaultHandler struct
	DefaultHandler struct {
		responseHandler gogrpcgatewayginresponse.Handler
	}
)

// NewDefaultHandler creates a new default response handler
//
// Parameters:
//
//   - responseHandler: The response handler (cannot be nil)
//
// Returns:
//
//   - *DefaultHandler: The default handler
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
//
// Parameters:
//
//   - ctx: The gin context
//   - err: The error to handle
func (d DefaultHandler) HandleError(
	ctx *gin.Context,
	err error,
) {
	// If there is no error, return
	if err == nil {
		return
	}

	// Create a new response based on the error
	response := goginresponse.NewErrorResponse(err)

	// Check if the error is a gRPC status error
	if _, ok := status.FromError(err); ok {
		d.responseHandler.HandleError(ctx, response)
	} else {
		ctx.JSON(http.StatusUnauthorized, response.Data())
	}
	ctx.Abort()
}
