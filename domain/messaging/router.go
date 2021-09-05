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
	"github.com/sonalys/letterme/domain/interfaces"
	"github.com/sonalys/letterme/domain/models"
)

// Router is responsible for consuming events from message broker and handling it.
// There are 2 types of incoming events: requests and responses.
//	- Requests: we need to process them using Service and send a message back.
//	- Response: we required it to other ms and need to get the response inside the context.
type Router struct {
	handlers         map[string]models.RouterHandler
	pendingResponses sync.Map
	*Configuration
	*Dependencies
}

type Configuration struct {
	ResponseTimeout time.Duration `json:"response_timeout"`
	ResponseChannel string        `json:"response_channel"`
}

type Dependencies struct {
	interfaces.Messaging
}

// NewRouter instantiates a new event router.
func NewRouter(ctx context.Context, c *Configuration, m *Dependencies) (*Router, error) {
	router := &Router{
		handlers:         map[string]models.RouterHandler{},
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
			pending := value.(*models.PendingResponse)
			if time.Since(pending.CreatedAt) > r.ResponseTimeout {
				pending.RespChan <- models.Response{
					Error: "response timed-out",
				}
				r.pendingResponses.Delete(key)
			}
			return true
		})

		time.Sleep(r.ResponseTimeout)
	}
}

func (r *Router) startConsumer() models.DeliveryHandler {
	return func(ctx context.Context, d models.Delivery) {
		// Incoming messages that should respond.
		if handler, ok := r.handlers[d.Type]; ok {
			m, err := handler(ctx, d)
			r.Publish(d.AppId, models.Message{
				ReplyTo: d.ReplyTo,
				Body:    models.NewResponse(m, err),
			})
			return
		}

		// Responses received that should be handled.
		if resp, ok := r.pendingResponses.Load(d.ReplyTo); ok {
			d.Acknowledger.Ack(d.DeliveryTag, true)

			msg := new(models.Response)
			if err := d.GetBody(msg); err != nil {
				resp.(*models.PendingResponse).RespChan <- models.Response{
					Error: err.Error(),
				}
			} else {
				resp.(*models.PendingResponse).RespChan <- models.Response{
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
func (r *Router) Communicate(queue string, m models.Message, dst interface{}) error {
	ch := make(chan models.Response, 1)
	defer close(ch)

	m.AppId = r.ResponseChannel
	m.ReplyTo = uuid.New().String()

	if err := r.Publish(queue, m); err != nil {
		return err
	}

	r.pendingResponses.Store(m.ReplyTo, &models.PendingResponse{
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

func (r *Router) AddHandler(eventType string, handler models.RouterHandler) {
	r.handlers[eventType] = handler
}
