package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
)

// DeliveryHandler is a handler func capable of managing deliveries.
type DeliveryHandler func(ctx context.Context, d Delivery)

type Table map[string]interface{}

type Response struct {
	Message Delivery
	Error   error
}

type PendingResponse struct {
	RespChan  chan<- Response
	CreatedAt time.Time
}

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
	Body interface{}
}

type Channel interface {
	Ack(tag uint64, ackBeforeReceiving bool) error
	Nack(tag uint64, ackBeforeReceiving bool, requeue bool) error
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
	body        []byte
}

func (d Delivery) GetBody(dst interface{}) error {
	return json.Unmarshal(d.body, dst)
}

func NewDelivery(d amqp.Delivery) Delivery {
	return Delivery{
		Acknowledger: d.Acknowledger,
		Message: Message{
			Headers:         Table(d.Headers),
			ContentType:     d.ContentType,
			ContentEncoding: d.ContentEncoding,
			DeliveryMode:    d.DeliveryMode,
			Priority:        d.Priority,
			CorrelationId:   d.CorrelationId,
			ReplyTo:         d.ReplyTo,
			Expiration:      d.Expiration,
			MessageId:       d.MessageId,
			Timestamp:       d.Timestamp,
			Type:            d.Type,
			UserId:          d.UserId,
			AppId:           d.AppId,
		},
		ConsumerTag:  d.ConsumerTag,
		MessageCount: d.MessageCount,
		DeliveryTag:  d.DeliveryTag,
		Redelivered:  d.Redelivered,
		Exchange:     d.Exchange,
		RoutingKey:   d.RoutingKey,
		body:         d.Body,
	}
}
