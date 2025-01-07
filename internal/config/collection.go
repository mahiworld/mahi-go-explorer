package config

// Collection is the collection names
type Collection struct {
	UserCollection string
}

// CreateCollection creates a new collection
func CreateCollection() *Collection {
	return &Collection{
		UserCollection: "users",
	}
}
