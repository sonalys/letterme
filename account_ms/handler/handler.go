package handler

import (
	cInt "github.com/sonalys/letterme/account_ms/interfaces"
	"github.com/sonalys/letterme/domain/messaging"
)

type Handler struct {
	cInt.Service
}

// RegisterHandlers registers all the eventType handlers for this ms.
func RegisterHandlers(r messaging.EventRouter, s cInt.Service) {
	handler := Handler{s}

	r.AddHandler(messaging.FetchEmailPublicInfo, handler.getAccountPublicInfo)
}
