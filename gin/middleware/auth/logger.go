package auth

import (
	gologger "github.com/ralvarezdev/go-logger"
	gologgerstatus "github.com/ralvarezdev/go-logger/status"
)

// Logger is the logger for the auth middleware
type Logger struct {
	logger gologger.Logger
}

// NewLogger is the logger for the auth middleware
func NewLogger(logger gologger.Logger) (*Logger, error) {
	// Check if the logger is nil
	if logger == nil {
		return nil, gologger.ErrNilLogger
	}

	return &Logger{logger: logger}, nil
}

// MethodNotSupported logs that the method is not supported
func (l *Logger) MethodNotSupported(method string) {
	l.logger.LogMessage(
		gologger.NewLogMessage(
			"Method not supported",
			gologgerstatus.Warning,
			method,
		),
	)
}

// BaseUriIsLongerThanFullPath logs that the base URI is longer than the full path
func (l *Logger) BaseUriIsLongerThanFullPath(fullPath string) {
	l.logger.LogMessage(
		gologger.NewLogMessage(
			"Base URI is longer than full path",
			gologgerstatus.Warning,
			fullPath,
		),
	)
}

// MissingGRPCMethod logs a MissingGRPCMethodError
func (l *Logger) MissingGRPCMethod(fullPath string) {
	l.logger.LogMessage(
		gologger.NewLogMessage(
			"Missing gRPC method",
			gologgerstatus.Warning,
			fullPath,
		),
	)
}

// MissingGRPCInterceptions logs a MissingGRPCInterceptionsError
func (l *Logger) MissingGRPCInterceptions() {
	l.logger.LogError(
		gologger.NewLogError(
			"Missing gRPC interceptions",
			ErrNilGRPCInterceptions,
		),
	)
}
