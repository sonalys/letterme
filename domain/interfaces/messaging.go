package interfaces

import (
	"context"

	"github.com/sonalys/letterme/domain/models"
)

type Messaging interface {
	Close() error
	CreateQueue(name string) error
	Publish(queue string, m models.Message) error
	Consume(ctx context.Context, queue string, handler models.DeliveryHandler) error
}
