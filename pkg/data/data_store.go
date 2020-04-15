package data

import (
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const collection = "data"

type IDataStore interface {
	Get(order string, skip, take int, startDate, finishDate time.Time)
	Save(data entities.Data)
}

type Store struct {
	Database *mgo.Database
	logger   logging.Logger
}

func NewStore(database *mgo.Database, logger logging.Logger) *Store {
	return &Store{database, logger}
}

func (ds *Store) Get(order string, skip, take int, startDate, finishDate time.Time) ([]entities.Data, error) {
	var data []entities.Data

	err := ds.Database.C(collection).Find(bson.M{
		"timestamp": bson.M{
			"$gt": startDate,
			"$lt": finishDate,
		},
	}).Select(bson.M{
		"timestamp": 1,
		"payload":   1,
		"from":      1,
	}).Skip(skip).Sort(order).Limit(take).All(&data)
	if err != nil {
		ds.logger.Error(err)
	}

	if data == nil {
		data = []entities.Data{}
	}

	return data, err
}

func (ds *Store) Save(data entities.Data) error {
	err := ds.Database.C(collection).Insert(&data)
	if err != nil {
		ds.logger.Error(err)
	}
	return err
}
