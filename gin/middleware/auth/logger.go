package auth

import (
	gologgermode "github.com/ralvarezdev/go-logger/mode"
	gologgermodenamed "github.com/ralvarezdev/go-logger/mode/named"
)

// Logger is the logger for the auth middleware
type Logger struct {
	logger gologgermodenamed.Logger
}

// NewLogger is the logger for the auth middleware
func NewLogger(header string, modeLogger gologgermode.Logger) (*Logger, error) {
	// Initialize the mode named logger
	namedLogger, err := gologgermodenamed.NewDefaultLogger(header, modeLogger)
	if err != nil {
		return nil, err
	}

	return &Logger{logger: namedLogger}, nil
}

// MethodNotSupported logs that the method is not supported
func (l *Logger) MethodNotSupported(method string) {
	l.logger.Warning(
		"method not supported",
		method,
	)
}

// BaseUriIsLongerThanFullPath logs that the base URI is longer than the full path
func (l *Logger) BaseUriIsLongerThanFullPath(fullPath string) {
	l.logger.Warning(
		"base uri is longer than full path",
		fullPath,
	)
}

// MissingGRPCMethod logs a MissingGRPCMethodError
func (l *Logger) MissingGRPCMethod(fullPath string) {
	l.logger.Warning(
		"missing grpc method",
		fullPath,
	)
}
