package store

import (
	"context"
	"log"
	"mahi-go-explorer/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// DbURL is the url of the database
	DbURL = config.GetFromEnv("DB_URL")
	// DbName is the name of the database
	DbName = config.GetFromEnv("DB_NAME")
)

// ConnectMongoDB connects to mongodb & returns the database
func ConnectMongoDB() (*mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(DbURL)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}

	//check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	//create indexes
	indices := []CollectionIndex{
		{
			Collection: *client.Database(DbName).Collection("users"),
			IndexKeys:  bson.D{{Key: "email", Value: 1}},
			Unique:     true,
		},
	}

	for _, index := range indices {
		err := createIndex(index)
		if err != nil {
			log.Println("Error in creating index", err.Error())
		}
	}

	return client.Database(DbName), nil
}

// CollectionIndex defines collection index
type CollectionIndex struct {
	Collection mongo.Collection
	IndexKeys  bson.D
	Unique     bool
}

func createIndex(index CollectionIndex) error {
	_, err := index.Collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    index.IndexKeys,
		Options: options.Index().SetUnique(index.Unique),
	})
	return err
}
