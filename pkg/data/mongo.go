package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// CONNECTIONTIMEOUT is a constant to set the waiting time for a database connection
const connectionTimeout = 60 * time.Second

type MongoDB struct {
	Server   string
	Database string
}

func NewMongoDB(server string, database string) *MongoDB {
	return &MongoDB{server, database}
}

func (s *MongoDB) Connect() (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(s.Server)
	clientOptions.SetConnectTimeout(connectionTimeout)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(s.Database), nil
}
