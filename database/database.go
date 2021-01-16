package database

import (
	"context"
	"go-merchants/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB configuration to establish connection with the db and host
func ConnectDB(ctx context.Context, conf config.MongoConfiguration) *mongo.Database {
	connection := options.Client().ApplyURI(conf.Server)
	// establish connection
	client, err := mongo.Connect(ctx, connection)
	if err != nil {
		panic(err)
	}

	return client.Database(conf.Database)
}
