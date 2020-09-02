package main

import (
	"context"
	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func disconnectClient(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func setupDatabase(databaseURI string, databaseName string, logger logging.Logger) (*mongo.Database, context.Context, *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseURI))
	failOnError(err, "Failed to connect to MongoDB")
	return client.Database(databaseName, nil), ctx, client
}

func setupStore(newDatabase *mongo.Database, logger logging.Logger, intValue int32) data.Store {
	return data.NewStore(newDatabase, logger, intValue)
}
