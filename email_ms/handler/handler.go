package handler

import (
	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/email_ms/controller"
)

// RegisterHandlers registers all the eventType handlers for this ms.
func RegisterHandlers(r messaging.EventRouter, s *controller.Service) {
}
