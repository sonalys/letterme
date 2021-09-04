package mongo

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/interfaces"
	dModels "github.com/sonalys/letterme/domain/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConfigEnv is used to get configs from env.
const ConfigEnv = "LM_MONGO_SETTINGS"

// Client is used to interface with persistence definitions.
type Client struct {
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
func (m *Client) Wait() <-chan bool {
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

// NewClient creates a new mongo controller instance, validates given configurations
// and awaits for mongo to connect, blocking the thread.
func NewClient(ctx context.Context, c *Configuration) (*Client, error) {
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

	return &Client{
		ctx:    ctx,
		client: client.Database(c.DBName),
	}, nil
}

// CreateCollection creates a collection in mongo.
func (m *Client) CreateCollection(colName string, indexes []map[string]interface{}) (interfaces.Collection, error) {
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
func (m *Client) DeleteCollection(ctx context.Context, colName string) error {
	if err := m.client.Collection(colName).Drop(ctx); err != nil {
		return newCollectionOperationError("delete", colName, err)
	}
	return nil
}

// GetCollection returns collection from mongo.
func (m *Client) GetCollection(colName string) interfaces.Collection {
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
func (c *Collection) Create(ctx context.Context, documents ...interface{}) ([]dModels.DatabaseID, error) {
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
