package context

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtginctx "github.com/ralvarezdev/go-jwt/gin/context"
	gojwtgrpc "github.com/ralvarezdev/go-jwt/grpc"
	"google.golang.org/grpc/metadata"
)

// GetOutgoingCtx returns a context with the raw token
//
// Parameters:
//
//   - ctx: the gin context
//
// Returns:
//
//   - context.Context: the gRPC context with the token in the metadata
//   - error: an error if any occurred during the process
func GetOutgoingCtx(ctx *gin.Context) (context.Context, error) {
	// Get the raw token from the context
	token, err := gojwtginctx.GetCtxToken(ctx)
	if err != nil {
		// Check if the token is missing
		if errors.Is(err, gojwt.ErrMissingTokenInContext) {
			return context.Background(), nil
		}
		return nil, err
	}

	// Append the token to the gRPC context
	grpcCtx := metadata.AppendToOutgoingContext(
		context.Background(),
		gojwtgrpc.AuthorizationMetadataKey,
		gojwt.BearerPrefix+" "+token,
	)

	return grpcCtx, nil
}
