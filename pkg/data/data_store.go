package data

import (
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const collection = "data"

type IDataStore interface {
	Get(query *entities.Query)
	Save(data entities.Data)
	Delete(deviceID string)
}

type Store struct {
	Database *mgo.Database
	logger   logging.Logger
}

func NewStore(database *mgo.Database, logger logging.Logger) *Store {
	return &Store{database, logger}
}

func (ds *Store) Get(query *entities.Query) ([]entities.Data, error) {
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

func (ds *Store) Save(data entities.Data) error {
	err := ds.Database.C(collection).Insert(&data)
	if err != nil {
		ds.logger.Error(err)
		return err
	}
	return nil
}

func (ds *Store) Delete(deviceID string) error {
	if deviceID == "" {
		return ds.removeAll(nil)
	}
	return ds.removeAll(bson.M{"from": deviceID})
}

func (ds *Store) removeAll(query interface{}) error {
	_, err := ds.Database.C(collection).RemoveAll(query)
	if err != nil {
		ds.logger.Error(err)
		return err
	}
	return nil
}

func (ds *Store) getFindQuery(query *entities.Query) (bson.M, error) {
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
