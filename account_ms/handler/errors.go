package handler

import "errors"

var (
	errDecode   = errors.New("failed to decode event")
	errInternal = errors.New("failed to process request")
)
