package interfaces

import "github.com/sonalys/letterme/domain/models"

type Router interface {
	WaitResponse(id string) <-chan models.Response
	AddHandler(eventType string, handler models.DeliveryHandler)
}
