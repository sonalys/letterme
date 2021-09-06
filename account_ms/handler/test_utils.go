package handler

import (
	"testing"

	"github.com/sonalys/letterme/account_ms/mocks"
)

type testHandler struct {
	handler *Handler
	mock    *mocks.Service
}

type testCase struct {
	name     string
	data     interface{}
	preRun   func(t *testing.T, mock *mocks.Service)
	expResp  interface{}
	expError error
}
