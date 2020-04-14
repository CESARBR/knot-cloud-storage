package interactor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	. "github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/CESARBR/knot-cloud-storage/pkg/users"
)

// DataInteractor represents the data layer interactor structure
type DataInteractor struct {
	UsersService users.Authenticator
	DataStore    *data.DataStore
	logger       logging.Logger
}

// NewDataInteractor creates a new data interactor instance
func NewDataInteractor(usersService users.Authenticator, dataStore *data.DataStore, logger logging.Logger) *DataInteractor {
	return &DataInteractor{usersService, dataStore, logger}
}

// GetAll retrieves all the data present in the storage
func (d *DataInteractor) GetAll(token string, order, skip, take int, startDate, finishDate time.Time) ([]Data, error) {
	err := d.Authenticate(token)
	if err != nil {
		return nil, err
	}

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

// GetByID retrieves data by it's ID from the storage, if present
func (d *DataInteractor) GetByID(token, id string, order, skip, take int, startDate, finishDate time.Time) ([]Data, error) {
	err := d.Authenticate(token)
	if err != nil {
		return nil, err
	}

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

// Save inserts data to the storage, if it doesn't exist already
func (d *DataInteractor) Save(token, id string, data []entities.Payload, timestamp time.Time) error {
	err := d.Authenticate(token)
	if err != nil {
		return err
	}

	for _, dt := range data {
		err = d.DataStore.Save(Data{From: id, Payload: dt, Timestamp: timestamp})
		if err != nil {
			return fmt.Errorf("error saving data %v: %w", data, err)
		}
	}

	return nil
}

// Authenticate verifies if the access token is valid.
func (d *DataInteractor) Authenticate(token string) error {
	err := d.UsersService.Authenticate(token)
	if err != nil {
		return err
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
