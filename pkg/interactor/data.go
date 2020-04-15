package interactor

import (
	"strconv"

	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
)

// GetAll retrieves all the data present in the storage
func (d *DataInteractor) GetAll(token string, query *entities.Query) ([]entities.Data, error) {
	err := d.Authenticate(token)
	if err != nil {
		return nil, err
	}

	selectOrder := "timestamp"
	if query.Order == -1 {
		selectOrder = "-timestamp"
	}

	data, err := d.DataStore.Get(selectOrder, query.Skip, query.Take, query.StartDate, query.FinishDate)
	if err != nil {
		d.logger.Error(err)
	}
	return data, err
}

// GetByID retrieves data by it's ID from the storage, if present
func (d *DataInteractor) GetByID(token string, query *entities.Query) ([]entities.Data, error) {
	err := d.Authenticate(token)
	if err != nil {
		return nil, err
	}

	selectOrder := "timestamp"
	if query.Order == -1 {
		selectOrder = "-timestamp"
	}

	s, err := strconv.ParseInt(query.ThingID, 10, 64)
	if err != nil {
		d.logger.Errorf("Error when trying to parse ID from string to int")
		return nil, err
	}

	data, err := d.DataStore.Get(selectOrder, query.Skip, query.Take, query.StartDate, query.FinishDate)
	data = filterDataBySensorID(data, int(s))
	return data, err
}

func filterDataBySensorID(data []entities.Data, sensorId int) []entities.Data {
	filteredData := make([]entities.Data, 0)
	for _, v := range data {
		if v.Payload.SensorId == sensorId {
			filteredData = append(filteredData, v)
		}
	}
	return filteredData
}
