package interfaces

import (
	"context"

	"github.com/sonalys/letterme/domain/messaging"
)

type Messaging interface {
	Close() error
	CreateQueue(name string) error
	Publish(queue string, m messaging.Message) error
	Consume(ctx context.Context, queue string, handler messaging.DeliveryHandler) error
}
