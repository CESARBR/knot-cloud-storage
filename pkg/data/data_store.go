package data

import (
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
	Database *mgo.Database
	logger   logging.Logger
}

// NewStore creates a new Store instance
func NewStore(database *mgo.Database, logger logging.Logger) Store {
	return &store{database, logger}
}

// Get returns data messages from the database
func (ds *store) Get(query *entities.Query) ([]entities.Data, error) {
	data := []entities.Data{}

	selectOrder := "timestamp"
	if query.Order == -1 {
		selectOrder = "-timestamp"
	}

	findQuery, err := ds.getFindQuery(query)
	if err != nil {
		return data, nil
	}

	err = ds.Database.C(collection).Find(findQuery).Select(bson.M{
		"timestamp": 1,
		"payload":   1,
		"from":      1,
	}).Skip(query.Skip).Sort(selectOrder).Limit(query.Take).All(&data)
	if err != nil {
		ds.logger.Error(err)
		return data, err
	}

	return data, nil
}

// Save stores data messages in the database
func (ds *store) Save(data entities.Data) error {
	err := ds.Database.C(collection).Insert(&data)
	if err != nil {
		ds.logger.Error(err)
		return err
	}
	return nil
}

// Delete removes data messages from the database
func (ds *store) Delete(deviceID string) error {
	if deviceID == "" {
		return ds.removeAll(nil)
	}
	return ds.removeAll(bson.M{"from": deviceID})
}

func (ds *store) removeAll(query interface{}) error {
	_, err := ds.Database.C(collection).RemoveAll(query)
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
