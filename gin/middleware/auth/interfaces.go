package auth

import (
	"github.com/gin-gonic/gin"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// Authenticator interface
	Authenticator interface {
		Authenticate(
			grpcMethod string,
			grpcInterceptions map[string]*gojwttoken.Token,
		) gin.HandlerFunc
	}
)
