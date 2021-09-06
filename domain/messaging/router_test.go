package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/sonalys/letterme/domain/messaging/rabbitmq"
	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RouterAllOk(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	const eventType = "mock-event"

	const timeout = 10 * time.Second
	ms1Config := &Configuration{ResponseTimeout: timeout, ResponseChannel: "ms-1"}
	ms2Config := &Configuration{ResponseTimeout: timeout, ResponseChannel: "ms-2"}

	rabbit, err := rabbitmq.NewClientFromEnv()
	require.NoError(t, err)

	router1, err := NewRouter(ctx, ms1Config, &Dependencies{Messaging: rabbit})
	require.NoError(t, err)
	router2, err := NewRouter(ctx, ms2Config, &Dependencies{Messaging: rabbit})
	require.NoError(t, err)

	defer func() {
		rabbit.DeleteQueue("ms-1")
		rabbit.DeleteQueue("ms-2")
	}()

	router2.AddHandler(eventType, func(ctx context.Context, d models.Delivery) (interface{}, error) {
		// Assertions
		out := new(string)
		err := d.GetBody(out)
		require.NoError(t, err)
		assert.Equal(t, "sender", *out)

		return "handler", nil
	})

	resp := new(string)
	err = router1.Communicate("ms-2", models.Message{
		Type: eventType,
		Body: "sender",
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, "handler", *resp)
}

func Test_RouterTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	const eventType = "mock-event"

	const timeout = 10 * time.Millisecond
	ms1Config := &Configuration{ResponseTimeout: timeout, ResponseChannel: "ms-1"}
	ms2Config := &Configuration{ResponseTimeout: timeout, ResponseChannel: "ms-2"}

	rabbit, err := rabbitmq.NewClientFromEnv()
	require.NoError(t, err)

	router1, err := NewRouter(ctx, ms1Config, &Dependencies{Messaging: rabbit})
	require.NoError(t, err)
	router2, err := NewRouter(ctx, ms2Config, &Dependencies{Messaging: rabbit})
	require.NoError(t, err)

	defer func() {
		cancel()
		rabbit.DeleteQueue("ms-1")
		rabbit.DeleteQueue("ms-2")
	}()

	router2.AddHandler(eventType, func(ctx context.Context, d models.Delivery) (interface{}, error) {
		time.Sleep(2 * timeout)
		return "handler", nil
	})

	resp := new(string)
	err = router1.Communicate("ms-2", models.Message{
		Type: eventType,
		Body: "sender",
	}, resp)
	assert.Error(t, err)
}