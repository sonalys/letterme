package interfaces

// Persistence encapsulates all required functions from a persistence integration.
// It should work for any choosen database that implements it: mongo, postgres, etc...
type Persistence interface {
	CreateCollection(colName string, indexes []string) error
	GetCollection(colName string) error
}

// Collection encapsulates all required functions expected from a persistence collection.
// a collection is a group of similar documents, all under the same indexes.
type Collection interface {
	First(filter, dst interface{}) error
	List(filter, dst interface{}) error
	Create(documents ...interface{}) error
	Update(documents ...interface{}) error
	Delete(documents ...interface{}) error
}
