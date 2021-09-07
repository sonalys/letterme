package messaging

import "context"

// Messenger is a controller capable of receiving and sending messages.
type Messenger interface {
	Close() error
	CreateQueue(name Queue) error
	Publish(queue Queue, m Message) error
	Consume(ctx context.Context, queue Queue, handler DeliveryHandler) error
}

// EventRouter is a controller capable of redirecting events to handlers.
type EventRouter interface {
	Communicate(queue Queue, m Message, dst interface{}) error
	AddHandler(eventType Event, handler RouterHandler)
}
