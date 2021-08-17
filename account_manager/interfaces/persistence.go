package interfaces

import (
	"context"

	"github.com/sonalys/letterme/domain/models"
)

// Persistence encapsulates all required functions from a persistence integration.
// It should work for any choosen database that implements it: mongo, postgres, etc...
type Persistence interface {
	Wait() <-chan bool
	CreateCollection(colName string, indexes []map[string]interface{}) (Collection, error)
	GetCollection(colName string) Collection
	DeleteCollection(ctx context.Context, colName string) error
}

// Collection encapsulates all required functions expected from a persistence collection.
// a collection is a group of similar documents, all under the same indexes.
type Collection interface {
	First(ctx context.Context, filter, dst interface{}) error
	List(ctx context.Context, filter, dst interface{}) error
	Create(ctx context.Context, documents ...interface{}) ([]models.DatabaseID, error)
	Update(ctx context.Context, filter, update interface{}) error
	Delete(ctx context.Context, filter interface{}) error
}
