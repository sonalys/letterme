package persistence

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/account_manager/interfaces"
	"github.com/sonalys/letterme/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoEnv is used to get configs from env.
const MongoEnv = "LM_MONGO_SETTINGS"

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

func (c Configuration) Validate() error {
	var errList []error
	if len(c.Hosts) == 0 {
		errList = append(errList, newEmptyFieldError("hosts"))
	}

	if c.SessionName == "" {
		errList = append(errList, newEmptyFieldError("session_name"))
	}

	if c.DBName == "" {
		errList = append(errList, newEmptyFieldError("db_name"))
	}

	if len(errList) > 0 {
		return newInvalidConfigurationError(errList)
	}
	return nil
}

const sleepTimeSeconds = 5

// Wait returns a chan that will return true when db is connected.
func (m *Mongo) Wait() <-chan bool {
	c := make(chan bool, 1)
	defer func() {
		for err := m.client.Client().Ping(m.ctx, nil); err != nil; {
			if m.ctx.Err() != nil {
				return
			}

			logrus.Infof("failed to connect to mongo: %s", err)
			time.Sleep(sleepTimeSeconds * time.Second)
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
		return nil, newInstanceError(err)
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, newInstanceError(err)
	}

	if err := client.Connect(ctx); err != nil {
		return nil, newConnectError(err)
	}

	return &Mongo{
		ctx:    ctx,
		client: client.Database(c.DBName),
	}, nil
}

// CreateCollection creates a collection in mongo.
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

// DeleteCollection deletes a collection from mongo.
func (m *Mongo) DeleteCollection(ctx context.Context, colName string) error {
	if err := m.client.Collection(colName).Drop(ctx); err != nil {
		return newCollectionOperationError("delete", colName, err)
	}
	return nil
}

// GetCollection returns collection from mongo.
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

// First Get the first document found, deserialize it inside dst.
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
func (c *Collection) Create(ctx context.Context, documents ...interface{}) ([]domain.DatabaseID, error) {
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
func (c *Collection) Update(ctx context.Context, filter, update interface{}) (int64, error) {
	res, err := c.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, newOperationError("updating", newCustomError(err))
	}
	return res.MatchedCount, nil
}

// Delete can remove one or more documents inside a collection,
// might return errors if one or more documents are not found.
func (c *Collection) Delete(ctx context.Context, filter interface{}) (int64, error) {
	res, err := c.Collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, newOperationError("deleting", newCustomError(err))
	}
	return res.DeletedCount, nil
}
