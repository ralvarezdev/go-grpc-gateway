package context

import (
	"context"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
)

// PrepareCtx prepares the context for the gRPC request
//
// Parameters:
//
//   - ctx: the gin context
//   - request: the request object to bind the JSON body to
//   - outgoingCtx: a function that takes the gin context and returns a context.Context and an error
//
// Returns:
//
//   - grpcCtx: the prepared gRPC context
//   - err: an error if any occurred during the process
func PrepareCtx(
	ctx *gin.Context,
	request interface{},
	outgoingCtx func(*gin.Context) (context.Context, error),
) (
	grpcCtx context.Context,
	err error,
) {
	// Bind the request
	if request != nil {
		err = ctx.ShouldBindJSON(request)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	return outgoingCtx(ctx)
}
