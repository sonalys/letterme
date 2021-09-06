package handler

import (
	cInt "github.com/sonalys/letterme/account_ms/interfaces"
	"github.com/sonalys/letterme/domain/interfaces"
	"github.com/sonalys/letterme/domain/messaging"
)

type Handler struct {
	cInt.Service
}

// RegisterHandlers registers all the eventType handlers for this ms.
func RegisterHandlers(r interfaces.Router, s cInt.Service) {
	handler := Handler{s}

	r.AddHandler(messaging.ECheckEmail, handler.verifyEmailExistence)
}
