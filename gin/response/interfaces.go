package response

import (
	"github.com/gin-gonic/gin"
	goginresponse "github.com/ralvarezdev/go-gin/response"
)

type (
	// Handler interface for handling the responses
	Handler interface {
		HandleSuccess(ctx *gin.Context, response *goginresponse.Response)
		HandleErrorProne(
			ctx *gin.Context,
			successResponse *goginresponse.Response,
			errorResponse *goginresponse.Response,
		)
		HandleError(ctx *gin.Context, response *goginresponse.Response)
	}
)
