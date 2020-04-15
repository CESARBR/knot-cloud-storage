package interactor

import (
	"strconv"

	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
)

func (d *DataInteractor) List(token string, query *entities.Query) ([]entities.Data, error) {
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

	if query.SensorID != "" {
		s, err := strconv.ParseInt(query.SensorID, 10, 64)
		if err != nil {
			d.logger.Errorf("Error when trying to parse ID from string to int")
			return nil, err
		}

		data = filterDataBySensorID(data, int(s))
	}

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
