package interactor

import (
	"strconv"
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	. "github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
)

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

func (d *DataInteractor) Save(data Data) error {
	err := d.DataStore.Save(data)
	if err != nil {
		d.logger.Error(err)
	}
	return err
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
