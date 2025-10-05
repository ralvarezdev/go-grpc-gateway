package route

import (
	"github.com/gin-gonic/gin"
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
)
