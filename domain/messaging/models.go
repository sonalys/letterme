package messaging

import (
	"context"
	"time"
)

type Table map[string]interface{}

type Message struct {
	// Application or exchange specific fields,
	// the headers exchange will inspect this field.
	Headers Table

	// Properties
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)
	Priority        uint8     // 0 to 9
	CorrelationId   string    // correlation identifier
	ReplyTo         string    // address to to reply to (ex: RPC)
	Expiration      string    // message expiration spec
	MessageId       string    // message identifier
	Timestamp       time.Time // message timestamp
	Type            string    // message type name
	UserId          string    // creating user id - ex: "guest"
	AppId           string    // creating application id

	// The application specific payload of the message
	Body []byte
}

type Delivery struct {
	Acknowledger Channel

	Message
	// Valid only with Channel.Consume
	ConsumerTag string
	// Valid only with Channel.Get
	MessageCount uint32

	DeliveryTag uint64
	Redelivered bool
	Exchange    string // basic.publish exchange
	RoutingKey  string // basic.publish routing key

}

// DeliveryHandler is a handler func capable of managing deliveries.
type DeliveryHandler func(ctx context.Context, d Delivery)
