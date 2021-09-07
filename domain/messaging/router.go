package messaging

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Router is responsible for consuming events from message broker and handling it.
// There are 2 types of incoming events: requests and responses.
//	- Requests: we need to process them using Service and send a message back.
//	- Response: we required it to other ms and need to get the response inside the context.
type Router struct {
	handlers         map[Event]RouterHandler
	pendingResponses sync.Map
	*Configuration
	*Dependencies
}

type Configuration struct {
	ResponseTimeout time.Duration `json:"response_timeout"`
	ResponseChannel Queue         `json:"response_channel"`
}

type Dependencies struct {
	Messenger
}

// NewRouter instantiates a new event router.
func NewRouter(ctx context.Context, c *Configuration, m *Dependencies) (*Router, error) {
	router := &Router{
		handlers:         map[Event]RouterHandler{},
		pendingResponses: sync.Map{},
		Configuration:    c,
		Dependencies:     m,
	}

	router.CreateQueue(c.ResponseChannel)

	err := router.Consume(ctx, c.ResponseChannel, router.startConsumer())
	if err != nil {
		return nil, err
	}
	go router.cleaningRoutine(ctx)

	return router, nil
}

func (r *Router) cleaningRoutine(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		r.pendingResponses.Range(func(key, value interface{}) bool {
			pending := value.(*PendingResponse)
			if time.Since(pending.CreatedAt) > r.ResponseTimeout {
				pending.RespChan <- Response{
					Error: "response timed-out",
				}
				r.pendingResponses.Delete(key)
			}
			return true
		})

		time.Sleep(r.ResponseTimeout)
	}
}

func (r *Router) startConsumer() DeliveryHandler {
	return func(ctx context.Context, d Delivery) {
		// Incoming messages that should respond.
		if handler, ok := r.handlers[d.Type]; ok {
			m, err := handler(ctx, d)
			r.Publish(d.AppId, Message{
				ReplyTo: d.ReplyTo,
				Body:    NewResponse(m, err),
			})
			return
		}

		// Responses received that should be handled.
		if resp, ok := r.pendingResponses.Load(d.ReplyTo); ok {
			d.Acknowledger.Ack(d.DeliveryTag, true)

			msg := new(Response)
			if err := d.GetBody(msg); err != nil {
				resp.(*PendingResponse).RespChan <- Response{
					Error: err.Error(),
				}
			} else {
				resp.(*PendingResponse).RespChan <- Response{
					Resp:  msg.Resp,
					Error: msg.Error,
				}
			}
			return
		}

		logrus.Errorf("failed to handle incoming event of type %s", d.Type)
		d.Acknowledger.Nack(d.DeliveryTag, true, false)
	}
}

// WaitResponse is used to retrieve an expected response from the queue.
func (r *Router) Communicate(queue Queue, m Message, dst interface{}) error {
	ch := make(chan Response, 1)
	defer close(ch)

	m.AppId = r.ResponseChannel
	m.ReplyTo = uuid.New().String()

	if err := r.Publish(queue, m); err != nil {
		return err
	}

	r.pendingResponses.Store(m.ReplyTo, &PendingResponse{
		RespChan:  ch,
		CreatedAt: time.Now(),
	})

	resp := <-ch
	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	valueDst := reflect.ValueOf(dst).Elem()
	if !valueDst.CanSet() {
		return fmt.Errorf("invalid dst type, %T is not writable", dst)
	}

	valueDst.Set(reflect.ValueOf(resp.Resp))
	return nil
}

func (r *Router) AddHandler(eventType Event, handler RouterHandler) {
	r.handlers[eventType] = handler
}
