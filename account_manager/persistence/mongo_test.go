package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

const mongo_env = "LM_MONGO_SETTINGS"

func loadFromEnv(key string, dst interface{}) error {
	if val, ok := os.LookupEnv(key); ok {
		if err := json.Unmarshal([]byte(val), dst); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("error loading config from env: '%s' not found", key)
}

func Test_Mongo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var cfg Configuration
	if err := loadFromEnv(mongo_env, &cfg); err != nil {
		require.Fail(t, err.Error())
	}

	mongo, err := NewMongo(ctx, &cfg)
	require.NoError(t, err, "should create without errors")
	require.NotNil(t, mongo, "mongo instance should exist")

	colName := "test-collection"
	col, err := mongo.CreateCollection(colName, []map[string]interface{}{
		{
			"name": 1,
		},
	})

	require.NoError(t, err, "should create collection")
	require.NotNil(t, col, "collection instance should exist")

	err = col.Delete(ctx, map[string]interface{}{})
	require.NoError(t, err, "should clear collection")

	defer func() {
		err = mongo.DeleteCollection(ctx, colName)
		require.NoError(t, err, "collection should be deleted")
	}()

	// Test
	t.Run("should create one document", func(t *testing.T) {
		doc := models.Account{
			OwnershipKey: "123",
		}

		ids, err := col.Create(ctx, doc)
		require.NoError(t, err, "document should be created")
		require.Len(t, ids, 1, "should return of created document")
	})

	t.Run("should create multiple documents", func(t *testing.T) {
		docs := []models.Account{
			{OwnershipKey: "123"},
			{OwnershipKey: "456"},
		}

		ids, err := col.Create(ctx, docs[0], docs[1])
		require.NoError(t, err, "document should be created")
		require.Len(t, ids, len(docs), "should return of created document")
	})

	t.Run("should get first document", func(t *testing.T) {
		var item models.Account
		err := col.First(ctx, models.Account{
			OwnershipKey: "456",
		}, &item)
		require.NoError(t, err)
		require.Equal(t, "456", item.OwnershipKey)
	})

	t.Run("should get list two document", func(t *testing.T) {
		var items []models.Account
		err := col.List(ctx, models.Account{
			OwnershipKey: "123",
		}, &items)
		require.NoError(t, err)
		require.Len(t, items, 2, "should return 2 items")
	})

	t.Run("should update one document", func(t *testing.T) {
		err := col.Update(ctx, models.Account{
			OwnershipKey: "456",
		}, map[string]interface{}{
			"$inc": models.Account{
				DeviceCount: 2,
			},
		})
		require.NoError(t, err)
		var item models.Account
		err = col.First(ctx, models.Account{
			OwnershipKey: "456",
		}, &item)
		require.NoError(t, err)
		require.Equal(t, uint8(2), item.DeviceCount)
		require.Equal(t, "456", item.OwnershipKey)
		require.NotEmpty(t, item.ID)
	})
}
