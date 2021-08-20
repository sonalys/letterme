package persistence

import (
	dModels "github.com/sonalys/letterme/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func convertMongoIDsToDatabaseIDs(ids []interface{}) ([]dModels.DatabaseID, error) {
	buf := make([]dModels.DatabaseID, 0, len(ids))
	for i := range ids {
		if id, ok := ids[i].(primitive.ObjectID); ok {
			buf = append(buf, dModels.DatabaseID(id.Hex()))
		} else {
			return nil, newCastError(primitive.ObjectID{}, dModels.DatabaseID(""))
		}
	}
	return buf, nil
}

func convertGenericIndexesToMongo(indexes []map[string]interface{}) []mongo.IndexModel {
	size := len(indexes)
	buf := make([]mongo.IndexModel, 0, size)
	for _, indexEntry := range indexes {
		keys := make([]bson.E, 0, len(indexEntry))
		for key, value := range indexEntry {
			keys = append(keys, bson.E{Key: key, Value: value})
		}
		buf = append(buf, mongo.IndexModel{
			Keys: keys,
		})
	}
	return buf
}
