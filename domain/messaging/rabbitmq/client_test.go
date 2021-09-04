package rabbitmq

import (
	"context"
	"sync"
	"testing"

	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func deleteQueue(c *Client, name string) {
	_ = c.DeleteQueue(name)
}

func Test_NewClient(t *testing.T) {
	client, err := NewClientFromEnv()
	require.NoError(t, err)
	require.NotNil(t, client)
}

func Test_CreateQueue(t *testing.T) {
	client, err := NewClientFromEnv()
	require.NoError(t, err)

	defer deleteQueue(client, "test")

	err = client.CreateQueue("test")
	require.NoError(t, err, "should create a queue")
	err = client.CreateQueue("test")
	require.NoError(t, err, "recreating existent queue shouldn't give error")
}

func Test_DeleteQueue(t *testing.T) {
	ctx := context.Background()
	client, err := NewClientFromEnv()
	require.NoError(t, err)

	err = client.CreateQueue("test")
	require.NoError(t, err, "should create a queue")

	err = client.DeleteQueue("test")
	require.NoError(t, err, "recreating existent queue shouldn't give error")

	err = client.Consume(ctx, "test", func(ctx context.Context, d models.Delivery) {})
	require.Error(t, err, "should give error for deleted queue")
}

func Test_PublishConsume(t *testing.T) {
	ctx := context.Background()
	client, err := NewClientFromEnv()
	require.NoError(t, err)

	queueName := "test"

	defer deleteQueue(client, queueName)

	err = client.CreateQueue(queueName)
	require.NoError(t, err, "should create a queue")

	msg := models.Message{
		Body: []byte{1, 2, 3},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	err = client.Consume(ctx, queueName, func(ctx context.Context, d models.Delivery) {
		defer wg.Done()
		resp := new([]byte)
		err := d.GetBody(resp)
		require.NoError(t, err)
		require.Equal(t, msg.Body, *resp)
	})
	require.NoError(t, err)

	err = client.Publish(queueName, msg)
	require.NoError(t, err, "publish in non-existent queue should create the queue")

	wg.Wait()
}
