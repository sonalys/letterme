package interfaces

import "github.com/sonalys/letterme/domain/models"

type Router interface {
	Communicate(queue string, m models.Message, dst interface{}) error
	AddHandler(eventType string, handler models.RouterHandler)
}
