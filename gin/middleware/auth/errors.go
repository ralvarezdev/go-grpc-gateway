package auth

import "errors"

var (
	ErrNilGRPCInterceptions = errors.New("grpc interceptions cannot be nil")
)
