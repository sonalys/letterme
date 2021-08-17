package persistence

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/account_manager/interfaces"
	"github.com/sonalys/letterme/domain/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MONGO_ENV is used to get configs from env.
const MONGO_ENV = "LM_MONGO_SETTINGS"

// Mongo is used to interface with persistence definitions.
type Mongo struct {
	ctx    context.Context
	client *mongo.Database
}

// Configuration is used to setup mongo connection.
type Configuration struct {
	Hosts       []string `json:"hosts"`
	SessionName string   `json:"session_name"`
	DBName      string   `json:"db_name"`
}

// Wait returns a chan that will return true when db is connected.
func (m *Mongo) Wait() <-chan bool {
	c := make(chan bool, 1)
	defer func() {
		for err := m.client.Client().Ping(m.ctx, nil); err != nil; {
			if m.ctx.Err() != nil {
				return
			}

			logrus.Infof("failed to connect to mongo: %s", err)
			time.Sleep(5 * time.Second)
		}
		c <- true
	}()

	return c
}

// NewMongo creates a new mongo controller instance, validates given configurations
// and awaits for mongo to connect, blocking the thread.
func NewMongo(ctx context.Context, c *Configuration) (*Mongo, error) {
	opts := options.Client().
		SetAppName(c.SessionName).
		SetHosts(c.Hosts)

	if err := opts.Validate(); err != nil {
		return nil, newInvalidConfigurationError(err)
	}

	if client, err := mongo.NewClient(opts); err == nil {
		if err := client.Connect(ctx); err != nil {
			return nil, newConnectError(err)
		}
		return &Mongo{
			ctx:    ctx,
			client: client.Database(c.DBName),
		}, nil
	} else {
		return nil, newInstanceError(err)
	}
}

func (m *Mongo) CreateCollection(colName string, indexes []map[string]interface{}) (interfaces.Collection, error) {
	createOpts := options.Collection()
	col := m.client.Collection(colName, createOpts)
	mongoIndexes := convertGenericIndexesToMongo(indexes)
	if _, err := col.Indexes().CreateMany(m.ctx, mongoIndexes, options.CreateIndexes()); err != nil {
		return nil, newCollectionOperationError("index", colName, err)
	}
	return &Collection{
		Collection: col,
	}, nil
}

func (m *Mongo) DeleteCollection(ctx context.Context, colName string) error {
	if err := m.client.Collection(colName).Drop(ctx); err != nil {
		return newCollectionOperationError("delete", colName, err)
	}
	return nil
}

func (m *Mongo) GetCollection(colName string) interfaces.Collection {
	createOpts := options.Collection()
	return &Collection{
		Collection: m.client.Collection(colName, createOpts),
	}
}

// Collection is mongo adapter for interfacing with persistence collections.
type Collection struct {
	*mongo.Collection
}

// Get the first document found, deserialize it inside dst.
// dst must be a struct otherwise will panic.
func (c *Collection) First(ctx context.Context, filter, dst interface{}) error {
	cur := c.Collection.FindOne(ctx, filter)
	err := cur.Err()

	switch err {
	case nil:
		if err := cur.Decode(dst); err != nil {
			return newOperationError("finding", newDecodeError(err))
		}
		return nil
	case mongo.ErrNoDocuments:
		return newOperationError("finding", newNotFoundError())
	default:
		return newOperationError("finding", newCustomError(err))
	}
}

// List returns all documents matched inside dst.
// dst must be a slice otherwise will panic.
func (c *Collection) List(ctx context.Context, filter, dst interface{}) error {
	cur, err := c.Collection.Find(ctx, filter)

	switch err {
	case nil:
		if err := cur.All(ctx, dst); err != nil {
			return newOperationError("listing", newDecodeError(err))
		}
		return nil
	case mongo.ErrNoDocuments:
		return newOperationError("listing", newNotFoundError())
	default:
		return newOperationError("listing", newCustomError(err))
	}
}

// Create can create one or multiple documents inside a collection,
// might get errors from invalid conversion, or empty documents list.
func (c *Collection) Create(ctx context.Context, documents ...interface{}) ([]models.DatabaseID, error) {
	cur, err := c.Collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, newOperationError("creating", newCustomError(err))
	}

	ids, err := convertMongoIDsToDatabaseIDs(cur.InsertedIDs)
	if err != nil {
		return nil, newOperationError("creating", newCustomError(err))
	}

	return ids, nil
}

// Update can update existent documents inside a collection,
// might return errors if one or more documents are not found.
func (c *Collection) Update(ctx context.Context, filter, update interface{}) error {
	_, err := c.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return newOperationError("updating", newCustomError(err))
	}
	return nil
}

// Delete can remove one or more documents inside a collection,
// might return errors if one or more documents are not found.
func (c *Collection) Delete(ctx context.Context, filter interface{}) error {
	_, err := c.Collection.DeleteMany(ctx, filter)
	if err != nil {
		return newOperationError("deleting", newCustomError(err))
	}
	return nil
}
