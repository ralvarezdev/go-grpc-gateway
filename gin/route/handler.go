package route

import (
	"github.com/gin-gonic/gin"
	gogrpcgatewayginauth "github.com/ralvarezdev/go-grpc-gateway/gin/middleware/auth"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// DefaultHandler struct
	DefaultHandler struct {
		authentication    gogrpcgatewayginauth.Authenticator
		grpcInterceptions map[string]*gojwttoken.Token
	}
)

// NewDefaultHandler creates a new default response handler
//
// Parameters:
//
//   - authentication: The authentication middleware (cannot be nil)
//   - grpcInterceptions: The gRPC interceptions map (cannot be nil)
//
// Returns:
//
//   - *DefaultHandler: The default handler
func NewDefaultHandler(
	authentication gogrpcgatewayginauth.Authenticator,
	grpcInterceptions map[string]*gojwttoken.Token,
) *DefaultHandler {
	return &DefaultHandler{
		authentication:    authentication,
		grpcInterceptions: grpcInterceptions,
	}
}

// New creates an authenticated endpoint if there is the access token or the refresh token required
//
// Parameters:
//
//   - route: The HTTP route
//   - grpcMethod: The gRPC method to authenticate
//   - handler: The handler function
//
// Returns:
//
//   - string: The HTTP route
//   - gin.HandlerFunc: The authentication middleware function (if required)
//   - gin.HandlerFunc: The handler function
func (d DefaultHandler) New(
	route, grpcMethod string,
	handler gin.HandlerFunc,
) (
	string,
	gin.HandlerFunc,
	gin.HandlerFunc,
) {
	// Create the endpoint
	return route, d.authentication.Authenticate(
		grpcMethod,
		d.grpcInterceptions,
	), handler
}
