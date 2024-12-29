package route

import (
	"github.com/gin-gonic/gin"
	gogrpcgatewayginauth "github.com/ralvarezdev/go-grpc-gateway/gin/middleware/auth"
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
)

type (
	// Handler interface
	Handler interface {
		New(route, grpcMethod string, handler gin.HandlerFunc) (
			string,
			gin.HandlerFunc,
			gin.HandlerFunc,
		)
	}

	// DefaultHandler struct
	DefaultHandler struct {
		authentication    gogrpcgatewayginauth.Authenticator
		grpcInterceptions *map[string]gojwtinterception.Interception
	}
)

// NewDefaultHandler creates a new default response handler
func NewDefaultHandler(
	authentication gogrpcgatewayginauth.Authenticator,
	grpcInterceptions *map[string]gojwtinterception.Interception,
) *DefaultHandler {
	return &DefaultHandler{
		authentication:    authentication,
		grpcInterceptions: grpcInterceptions,
	}
}

// New creates an authenticated endpoint if there is the access token or the refresh token required
func (d *DefaultHandler) New(
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
