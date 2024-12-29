package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gongin "github.com/ralvarezdev/go-gin"
	gongintypes "github.com/ralvarezdev/go-gin/response"
	gogrpcstauts "github.com/ralvarezdev/go-grpc/status"
	"google.golang.org/grpc/codes"
	"net/http"
)

type (
	// Handler interface
	Handler interface {
		HandlePrepareCtxError(ctx *gin.Context, err error)
		HandleResponse(ctx *gin.Context, response interface{}, err error)
		HandleErrorResponse(ctx *gin.Context, err error)
	}

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

// HandleResponse handles the response from the gRPC server
func (d *DefaultHandler) HandleResponse(
	ctx *gin.Context,
	code int,
	response interface{},
	err error,
) {
	// Check if the error is nil
	if err == nil {
		ctx.JSON(code, response)
		return
	}

	// Handle the error response
	d.HandleErrorResponse(ctx, err)
}

// HandleErrorResponse handles the error response from the gRPC server
func (d *DefaultHandler) HandleErrorResponse(ctx *gin.Context, err error) {
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
