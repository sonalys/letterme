package mongo

import (
	"testing"

	dModels "github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_convertMongoIDsToDatabaseIDs(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()
		ids := []interface{}{id1, id2}
		got, err := convertMongoIDsToDatabaseIDs(ids)
		require.NoError(t, err)

		expected := []dModels.DatabaseID{dModels.DatabaseID(id1.Hex()), dModels.DatabaseID(id2.Hex())}
		require.Equal(t, expected, got)
	})

	t.Run("invalid id", func(t *testing.T) {
		id1 := primitive.NewObjectID()
		id2 := "123"
		ids := []interface{}{id1, id2}
		got, err := convertMongoIDsToDatabaseIDs(ids)
		require.Error(t, err, "should give err for invalid id")
		require.Nil(t, got)
	})
}

func Test_convertGenericIndexesToMongo(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		indexes := []map[string]interface{}{
			{
				"myKey1": -1,
			},
			{
				"myKey1": 1,
				"myKey2": 1,
			},
		}

		expected := []mongo.IndexModel{
			{
				Keys: []bson.E{
					{
						Key:   "myKey1",
						Value: -1,
					},
				},
			},
			{
				Keys: []bson.E{
					{
						Key:   "myKey1",
						Value: 1,
					},
					{
						Key:   "myKey2",
						Value: 1,
					},
				},
			},
		}

		got := convertGenericIndexesToMongo(indexes)

		for i := range expected {
			require.EqualValues(t, expected[i], got[i])
		}
	})
}
