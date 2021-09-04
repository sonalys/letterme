package messaging

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/interfaces"
	"github.com/sonalys/letterme/domain/models"
)

// Router is responsible for consuming events from message broker and handling it.
// There are 2 types of incoming events: requests and responses.
//	- Requests: we need to process them using Service and send a message back.
//	- Response: we required it to other ms and need to get the response inside the context.
type Router struct {
	handlers         map[string]models.DeliveryHandler
	pendingResponses sync.Map
	*Configuration
	*Dependencies
}

type Configuration struct {
	ResponseTimeout time.Duration `json:"response_timeout"`
}

type Dependencies struct {
	interfaces.Messaging
}

// NewRouter instantiates a new event router.
func NewRouter(ctx context.Context, c *Configuration, m *Dependencies) *Router {
	router := &Router{
		handlers:         map[string]models.DeliveryHandler{},
		pendingResponses: sync.Map{},
		Configuration:    c,
		Dependencies:     m,
	}
	go router.Consume(ctx, QEmailMS, router.startConsumer())
	go router.cleaningRoutine(ctx)
	return router
}

func (r *Router) cleaningRoutine(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		r.pendingResponses.Range(func(key, value interface{}) bool {
			pending := value.(*models.PendingResponse)
			if time.Since(pending.CreatedAt) > r.ResponseTimeout {
				r.pendingResponses.Delete(key)
			}
			return true
		})

		time.Sleep(r.ResponseTimeout)
	}
}

func (r *Router) startConsumer() models.DeliveryHandler {
	return func(ctx context.Context, d models.Delivery) {
		if handler, ok := r.handlers[d.Type]; ok {
			handler(ctx, d)
			return
		}

		if resp, ok := r.pendingResponses.Load(d.ReplyTo); ok {
			d.Acknowledger.Ack(d.DeliveryTag, true)
			resp.(*models.PendingResponse).RespChan <- models.Response{
				Message: d,
			}
			return
		}

		logrus.Errorf("failed to handle incoming event of type %s", d.Type)
		d.Acknowledger.Nack(d.DeliveryTag, true, false)
	}
}

// WaitResponse is used to retrieve an expected response from the queue.
func (r *Router) WaitResponse(id string) <-chan models.Response {
	ch := make(chan models.Response, 1)
	r.pendingResponses.Store(id, &models.PendingResponse{
		RespChan:  ch,
		CreatedAt: time.Now(),
	})
	return ch
}

func (r *Router) AddHandler(eventType string, handler models.DeliveryHandler) {
	r.handlers[eventType] = handler
}
