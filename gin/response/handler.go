package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gongin "github.com/ralvarezdev/go-gin"
	goginresponse "github.com/ralvarezdev/go-gin/response"
	gongintypes "github.com/ralvarezdev/go-gin/response"
	gogrpcstauts "github.com/ralvarezdev/go-grpc/status"
	"google.golang.org/grpc/codes"
	"net/http"
)

type (
	// DefaultHandler struct
	DefaultHandler struct {
		mode *goflagsmode.Flag
	}
)

// NewDefaultHandler creates a new default request handler
func NewDefaultHandler(mode *goflagsmode.Flag) (*DefaultHandler, error) {
	// Check if the flag mode is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	return &DefaultHandler{mode: mode}, nil
}

// HandleSuccess handles the success response
func (d *DefaultHandler) HandleSuccess(
	ctx *gin.Context,
	response *goginresponse.Response,
) {
	if response != nil && response.Code != nil {
		ctx.JSON(*response.Code, response.Data)
	} else {
		goginresponse.SendInternalServerError(ctx)
	}
}

// HandleErrorProne handles the response that may contain an error
func (d *DefaultHandler) HandleErrorProne(
	ctx *gin.Context,
	successResponse *goginresponse.Response,
	errorResponse *goginresponse.Response,
) {
	// Check if the error response is nil
	if errorResponse != nil {
		d.HandleError(ctx, errorResponse)
		return
	}

	// Handle the success response
	d.HandleSuccess(ctx, successResponse)
}

// HandleError handles the error response
func (d *DefaultHandler) HandleError(
	ctx *gin.Context,
	response *goginresponse.Response,
) {
	// Check if the response is nil or if the response code is not nil
	if response == nil {
		goginresponse.SendInternalServerError(ctx)
		return
	} else if response.Code != nil {
		ctx.JSON(*response.Code, response.Data)
		ctx.Abort()
		return
	}

	// Get the error from the response data
	err, ok := response.Data.(error)
	if !ok {
		goginresponse.SendInternalServerError(ctx)
		return
	}

	// Extract the gRPC code and error from the status
	extractedCode, extractedErr := gogrpcstauts.ExtractErrorFromStatus(
		d.mode,
		err,
	)

	// Check the extracted code and error
	switch extractedCode {
	case codes.AlreadyExists:
		ctx.JSON(
			http.StatusConflict,
			gongintypes.NewErrorResponse(extractedErr),
		)
	case codes.NotFound:
		ctx.JSON(
			http.StatusNotFound,
			gongintypes.NewErrorResponse(extractedErr),
		)
	case codes.InvalidArgument:
		ctx.JSON(
			http.StatusBadRequest,
			gongintypes.NewErrorResponse(extractedErr),
		)
	case codes.PermissionDenied:
		if d.mode == nil || d.mode.IsProd() {
			ctx.JSON(
				http.StatusForbidden,
				gongintypes.NewErrorResponse(errors.New(gongin.Unauthorized)),
			)
		}
		ctx.JSON(
			http.StatusForbidden,
			gongintypes.NewErrorResponse(extractedErr),
		)
	case codes.Unauthenticated:
		if d.mode == nil || d.mode.IsProd() {
			ctx.JSON(
				http.StatusUnauthorized,
				gongintypes.NewErrorResponse(gongin.Unauthenticated),
			)
		}
		ctx.JSON(
			http.StatusUnauthorized,
			gongintypes.NewErrorResponse(extractedErr),
		)
	case codes.Unimplemented:
		if d.mode == nil || d.mode.IsProd() {
			ctx.JSON(
				http.StatusNotImplemented,
				gongintypes.NewErrorResponse(gongin.InDevelopment),
			)
		}
		ctx.JSON(
			http.StatusNotImplemented,
			gongintypes.NewErrorResponse(extractedErr),
		)
	case codes.Unavailable:
		if d.mode == nil || d.mode.IsProd() {
			ctx.JSON(
				http.StatusServiceUnavailable,
				gongintypes.NewErrorResponse(errors.New(gongin.ServiceUnavailable)),
			)
		}
		ctx.JSON(
			http.StatusServiceUnavailable,
			gongintypes.NewErrorResponse(extractedErr),
		)
	default:
		if d.mode == nil || d.mode.IsProd() {
			ctx.JSON(
				http.StatusInternalServerError,
				gongintypes.NewErrorResponse(errors.New(gongin.InternalServerError)),
			)
		}
		ctx.JSON(
			http.StatusInternalServerError,
			gongintypes.NewErrorResponse(extractedErr),
		)
	}
}
