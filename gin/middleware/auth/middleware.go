package auth

import (
	"github.com/gin-gonic/gin"
	goginmiddlewareauth "github.com/ralvarezdev/go-gin/middleware/auth"
	goginresponse "github.com/ralvarezdev/go-gin/response"
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
	gologger "github.com/ralvarezdev/go-logger"
)

// Middleware struct
type Middleware struct {
	logger          *Logger
	authenticator   goginmiddlewareauth.Authenticator
	authenticateFns map[string]gin.HandlerFunc
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(
	logger *Logger,
	authenticator goginmiddlewareauth.Authenticator,
) (*Middleware, error) {
	// Check if either the logger or authenticator is nil
	if logger == nil {
		return nil, gologger.ErrNilLogger
	}
	if authenticator == nil {
		return nil, goginmiddlewareauth.ErrNilAuthenticator
	}

	return &Middleware{
		logger:        logger,
		authenticator: authenticator,
	}, nil
}

// Authenticate return the middleware function that authenticates the request
func (m *Middleware) Authenticate(
	grpcMethod string,
	grpcInterceptions *map[string]gojwtinterception.Interception,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if the gRPC interceptions is nil
		if grpcInterceptions == nil {
			if grpcInterceptions == nil {
				m.logger.MissingGRPCInterceptions()
			}
			goginresponse.SendInternalServerError(ctx)
			return
		}

		// Get the request URI and method
		requestURI := ctx.Request.RequestURI

		// Get the gRPC method interception
		interception, ok := (*grpcInterceptions)[grpcMethod]
		if !ok {
			m.logger.MissingGRPCMethod(requestURI)
			goginresponse.SendInternalServerError(ctx)
			return
		}

		// Check if there is None interception
		if interception == gojwtinterception.None {
			ctx.Next()
			return
		}

		// Check if the interception authentication function is already set
		fn, ok := m.authenticateFns[grpcMethod]
		if !ok {
			fn = m.authenticator.Authenticate(interception)
			m.authenticateFns[grpcMethod] = fn
		}
		fn(ctx)
	}
}
