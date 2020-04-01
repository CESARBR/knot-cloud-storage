package interactor

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	. "github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
)

// ErrTokenEmpty is returned for empty access tokens.
var ErrTokenEmpty = errors.New("no access token provided")

type DataInteractor struct {
	DataStore *data.DataStore
	logger    logging.Logger
}

func NewDataInteractor(dataStore *data.DataStore, logger logging.Logger) *DataInteractor {
	return &DataInteractor{dataStore, logger}
}

func (d *DataInteractor) GetAll(order, skip, take int, startDate, finishDate time.Time) ([]Data, error) {
	selectOrder := "timestamp"
	if order == -1 {
		selectOrder = "-timestamp"
	}

	data, err := d.DataStore.Get(selectOrder, skip, take, startDate, finishDate)
	if err != nil {
		d.logger.Error(err)
	}
	return data, err
}

func (d *DataInteractor) GetByID(id string, order, skip, take int, startDate, finishDate time.Time) ([]Data, error) {
	selectOrder := "timestamp"
	if order == -1 {
		selectOrder = "-timestamp"
	}

	s, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		d.logger.Errorf("Error when trying to parse ID from string to int")
		return nil, err
	}

	data, err := d.DataStore.Get(selectOrder, skip, take, startDate, finishDate)
	data = filterDataBySensorID(data, int(s))
	return data, err
}

func (d *DataInteractor) Save(id string, data []entities.Payload, timestamp time.Time) error {
	for _, dt := range data {
		err := d.DataStore.Save(Data{From: id, Payload: dt, Timestamp: timestamp})
		if err != nil {
			return fmt.Errorf("error saving data %v: %w", data, err)
		}
	}

	return nil
}

// Authenticate verifies if the access token is valid.
func (d *DataInteractor) Authenticate(token string) error {
	if token == "" {
		return ErrTokenEmpty
	}
	return nil
}

func filterDataBySensorID(data []Data, sensorId int) []Data {
	filteredData := make([]Data, 0)
	for _, v := range data {
		if v.Payload.SensorId == sensorId {
			filteredData = append(filteredData, v)
		}
	}
	return filteredData
}
