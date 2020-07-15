package data

import (
	"context"
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collection = "data"

// Store represents the interface to data related operations
type Store interface {
	Get(query *entities.Query) ([]entities.Data, error)
	Save(data entities.Data) error
	Delete(deviceID string) error
}

// store represents the data capabilities implementation
type store struct {
	Database *mongo.Database
	logger   logging.Logger
}

// NewStore creates a new Store instance
func NewStore(database *mongo.Database, logger logging.Logger, expTime int32) Store {
	err := setupDataExpiration(database, expTime)
	if err != nil {
		logger.Infof("fail to set data expiration time: %w", err)
	}

	return &store{database, logger}
}

func setupDataExpiration(db *mongo.Database, time int32) error {
	if time == 0 {
		return nil
	}

	_, err := db.Collection(collection).Indexes().DropAll(context.TODO())
	if err != nil {
		return err
	}

	_, err = db.Collection(collection).Indexes().CreateOne(context.TODO(), getTimestampIndex(time))
	if err != nil {
		return err
	}

	return nil
}

func getTimestampIndex(time int32) mongo.IndexModel {
	return mongo.IndexModel{
		Keys: bson.M{
			"timestamp": 1,
		},
		Options: &options.IndexOptions{
			ExpireAfterSeconds: &time,
		},
	}
}

// Get returns data messages from the database
func (ds *store) Get(query *entities.Query) ([]entities.Data, error) {
	data := []entities.Data{}

	findQuery, err := ds.getFindQuery(query)
	if err != nil {
		return data, nil
	}

	options := ds.getFindOptions(query)
	cur, err := ds.Database.Collection(collection).Find(context.TODO(), findQuery, options)
	if err != nil {
		ds.logger.Error(err)
		return data, err
	}

	for cur.Next(context.TODO()) {
		var decodedData entities.Data
		err = cur.Decode(&decodedData)
		if err != nil {
			return data, err
		}

		data = append(data, decodedData)
	}
	return data, nil
}

// Save stores data messages in the database
func (ds *store) Save(data entities.Data) error {
	data.Timestamp = time.Now()
	_, err := ds.Database.Collection(collection).InsertOne(context.TODO(), &data)
	if err != nil {
		ds.logger.Error(err)
		return err
	}
	return nil
}

// Delete removes data messages from the database
func (ds *store) Delete(deviceID string) error {
	if deviceID == "" {
		return ds.removeAll(bson.M{})
	}
	return ds.removeAll(bson.M{"from": deviceID})
}

func (ds *store) removeAll(query interface{}) error {
	_, err := ds.Database.Collection(collection).DeleteMany(context.TODO(), query)
	if err != nil {
		ds.logger.Error(err)
		return err
	}
	return nil
}

func (ds *store) getFindQuery(query *entities.Query) (bson.M, error) {
	b := bson.M{
		"timestamp": bson.M{
			"$gt": query.StartDate,
			"$lt": query.FinishDate,
		},
	}

	if query.ThingID != "" {
		b["from"] = query.ThingID
	}

	return b, nil
}

func (ds *store) getFindOptions(query *entities.Query) *options.FindOptions {
	options := options.Find()
	options.SetProjection(
		bson.M{
			"timestamp": 1,
			"payload":   1,
			"from":      1,
		},
	)
	options.SetSkip(query.Skip)
	options.SetLimit(query.Take)
	options.SetSort(bson.M{"timestamp": query.Order})

	return options
}
