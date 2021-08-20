package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/sonalys/letterme/account_manager/utils"
	dModels "github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func Test_Mongo(t *testing.T) {
	ctx := context.Background()
	var cfg Configuration
	if err := utils.LoadFromEnv(MongoEnv, &cfg); err != nil {
		require.Fail(t, err.Error())
	}

	mongo, err := NewMongo(ctx, &cfg)
	require.NoError(t, err, "should create without errors")
	require.NotNil(t, mongo, "mongo instance should exist")

	select {
	case <-mongo.Wait():
		break
	case <-time.After(5 * time.Second):
		require.Fail(t, "database connection timedout")
	}

	colName := "test-collection"
	col, err := mongo.CreateCollection(colName, []map[string]interface{}{
		{
			"name": 1,
		},
	})

	require.NoError(t, err, "should create collection")
	require.NotNil(t, col, "collection instance should exist")

	_, err = col.Delete(ctx, map[string]interface{}{})
	require.NoError(t, err, "should clear collection")

	defer func() {
		err = mongo.DeleteCollection(ctx, colName)
		require.NoError(t, err, "collection should be deleted")
	}()

	// Test
	t.Run("should create one document", func(t *testing.T) {
		doc := dModels.Account{
			OwnershipKey: dModels.OwnershipKey("123"),
		}

		ids, err := col.Create(ctx, doc)
		require.NoError(t, err, "document should be created")
		require.Len(t, ids, 1, "should return of created document")
	})

	t.Run("should create multiple documents", func(t *testing.T) {
		docs := []dModels.Account{
			{OwnershipKey: dModels.OwnershipKey("123")},
			{OwnershipKey: dModels.OwnershipKey("456")},
		}

		ids, err := col.Create(ctx, docs[0], docs[1])
		require.NoError(t, err, "document should be created")
		require.Len(t, ids, len(docs), "should return of created document")
	})

	t.Run("should get first document", func(t *testing.T) {
		var item dModels.Account
		err := col.First(ctx, dModels.Account{
			OwnershipKey: dModels.OwnershipKey("456"),
		}, &item)
		require.NoError(t, err)
		require.Equal(t, dModels.OwnershipKey("456"), item.OwnershipKey)
	})

	t.Run("should get list two document", func(t *testing.T) {
		var items []dModels.Account
		err := col.List(ctx, dModels.Account{
			OwnershipKey: dModels.OwnershipKey("123"),
		}, &items)
		require.NoError(t, err)
		require.Len(t, items, 2, "should return 2 items")
	})

	t.Run("should update one document", func(t *testing.T) {
		count, err := col.Update(ctx, dModels.Account{
			OwnershipKey: dModels.OwnershipKey("456"),
		}, map[string]interface{}{
			"$inc": dModels.Account{
				DeviceCount: 2,
			},
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), count, "should update exactly one document")
		var item dModels.Account
		err = col.First(ctx, dModels.Account{
			OwnershipKey: dModels.OwnershipKey("456"),
		}, &item)
		require.NoError(t, err)
		require.Equal(t, uint8(2), item.DeviceCount)
		require.Equal(t, dModels.OwnershipKey("456"), item.OwnershipKey)
		require.NotEmpty(t, item.ID)
	})
}
