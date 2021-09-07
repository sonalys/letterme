package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/utils"
	"github.com/streadway/amqp"
)

const ConfigEnv = "LM_RABBIQMQ_SETTINGS"

// Delivery is the encapsulation of the amqp.Delivery
type Delivery = amqp.Delivery

// Message is the encapsulation of the amqp.Message
type Message = amqp.Publishing

// Channel is the encapsulation of the amqp.Channel
type Channel = amqp.Channel

// Client is an adapter for the rabbiqMQ
type Client struct {
	adapter   *amqp.Connection
	session   *Channel
	consumers sync.Map
	*Configuration
}

type Configuration struct {
	Host       string `json:"host"`
	MaxRetries uint
}

// NewClient returns a new rabbitMQ client.
func NewClient(c *Configuration) (*Client, error) {
	adapter, err := amqp.Dial(c.Host)
	if err != nil {
		return nil, errors.Wrap(err, errDialing)
	}

	return &Client{
		adapter: adapter,
	}, nil
}

func NewClientFromEnv() (*Client, error) {
	config := new(Configuration)
	if err := utils.LoadFromEnv(ConfigEnv, config); err != nil {
		return nil, errors.Wrap(err, errDialing)
	}

	return NewClient(config)
}

func (c *Client) Close() error {
	if err := c.adapter.Close(); err != nil {
		return errors.Wrap(err, errClosing)
	}

	return nil
}

func (c *Client) revalidateSession() error {
	var ch *amqp.Channel
	var err error

	var retries uint
	for ch, err = c.adapter.Channel(); err != nil && retries <= c.MaxRetries; {
		if err != nil {
			retries++
			logrus.Error(errors.Wrap(err, errCreateQueue))
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

	c.session = ch
	return err
}

// getChannel gets a new connection with rabbitMQ.
func (c *Client) getChannel() (*Channel, error) {
	if c.session == nil {
		err := c.revalidateSession()
		if err != nil {
			return nil, err
		}
	}

	return c.session, nil
}

func (c *Client) DeleteQueue(name messaging.Queue) error {
	ch, err := c.getChannel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDelete(string(name), true, true, true)
	if err != nil {
		return err
	}

	return nil
}

// CreateQueue creates a new topic in which you can either publish or subscribe.
func (c *Client) CreateQueue(name messaging.Queue) error {
	ch, err := c.getChannel()
	if err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(
		string(name), false, false, false, false, nil,
	); err != nil {
		_ = c.revalidateSession()
		return err
	}

	return nil
}

// Publish sends a new message to the specified queue, the queue must already exist.
func (c *Client) Publish(queue messaging.Queue, m messaging.Message) error {
	ch, err := c.getChannel()
	if err != nil {
		return err
	}

	if err := ch.Publish("", string(queue), false, false, transformMessageToRabbit(m)); err != nil {
		_ = c.revalidateSession()
		return errors.Wrap(err, fmt.Sprintf(errPublish, queue))
	}

	return nil
}

func transformMessageToRabbit(m messaging.Message) amqp.Publishing {
	body, _ := json.Marshal(m.Body)
	return amqp.Publishing{
		Headers:         amqp.Table(m.Headers),
		ContentType:     m.ContentType,
		ContentEncoding: m.ContentEncoding,
		DeliveryMode:    m.DeliveryMode,
		Priority:        m.Priority,
		CorrelationId:   m.CorrelationId,
		ReplyTo:         m.ReplyTo,
		Expiration:      m.Expiration,
		MessageId:       m.MessageId,
		Timestamp:       m.Timestamp,
		Type:            string(m.Type),
		UserId:          m.UserId,
		AppId:           string(m.AppId),
		Body:            body,
	}
}

// Consume allows you to specify a handler for a given queue, the queue must already exist.
func (c *Client) Consume(ctx context.Context, queue messaging.Queue, handler messaging.DeliveryHandler) error {
	ch, err := c.getChannel()
	if err != nil {
		return err
	}

	recv, err := ch.Consume(string(queue), "", false, false, false, false, nil)
	if err != nil {
		_ = c.revalidateSession()
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = ch.Close()
				return
			case delivery := <-recv:
				handler(ctx, transformDeliveryFromRabbit(delivery))
			}
		}
	}()
	return nil
}

func transformDeliveryFromRabbit(d Delivery) messaging.Delivery {
	return messaging.NewDelivery(d)
}
