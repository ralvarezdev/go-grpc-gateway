package auth

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	goginmiddlewareauth "github.com/ralvarezdev/go-gin/middleware/auth"
	goginresponse "github.com/ralvarezdev/go-gin/response"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// Middleware struct
	Middleware struct {
		logger          *slog.Logger
		authenticator   goginmiddlewareauth.Authenticator
		authenticateFns map[string]gin.HandlerFunc
	}
)

// NewMiddleware creates a new authentication middleware
//
// Parameters:
//
//   - authenticator: The authenticator (cannot be nil)
//   - logger: The logger (optional, can be nil)
//
// Returns:
//
//   - *Middleware: The authentication middleware
func NewMiddleware(
	authenticator goginmiddlewareauth.Authenticator,
	logger *slog.Logger,
) (*Middleware, error) {
	// Check if either the authenticator is nil
	if authenticator == nil {
		return nil, goginmiddlewareauth.ErrNilAuthenticator
	}

	if logger != nil {
		logger = logger.With(slog.String("component", "gin_middleware_auth"))
	}

	return &Middleware{
		logger:        logger,
		authenticator: authenticator,
	}, nil
}

// Authenticate return the middleware function that authenticates the request
//
// Parameters:
//
//   - grpcMethod: The gRPC method to authenticate
//   - grpcInterceptions: The gRPC interceptions map
//
// Returns:
//
//   - gin.HandlerFunc: The middleware function
func (m Middleware) Authenticate(
	grpcMethod string,
	grpcInterceptions map[string]*gojwttoken.Token,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if the gRPC interceptions is nil
		if grpcInterceptions == nil {
			goginresponse.SendInternalServerError(ctx)
			return
		}

		// Get the request URI and method
		requestURI := ctx.Request.RequestURI

		// Get the gRPC method interception
		interception, ok := grpcInterceptions[grpcMethod]
		if !ok {
			if m.logger != nil {
				m.logger.Warn(
					"missing grpc method",
					slog.String("method", grpcMethod),
					slog.String("request_uri", requestURI),
				)
			}
			goginresponse.SendInternalServerError(ctx)
			return
		}

		// Check if there is None interception
		if interception == nil {
			ctx.Next()
			return
		}

		// Check if the interception authentication function is already set
		fn, ok := m.authenticateFns[grpcMethod]
		if !ok {
			fn = m.authenticator.Authenticate(*interception)
			m.authenticateFns[grpcMethod] = fn
		}
		fn(ctx)
	}
}
