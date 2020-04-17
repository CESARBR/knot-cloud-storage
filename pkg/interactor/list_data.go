package interactor

import (
	"errors"
	"fmt"
	"strconv"

	btEntities "github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
)

// List provides all the thing's data owned by an user
func (d *DataInteractor) List(token string, query *entities.Query) ([]entities.Data, error) {
	data, err := d.getDevicesData(token, query)
	if err != nil {
		return []entities.Data{}, err
	}

	data = d.filterDataBySensorID(data, query.SensorID)
	return data, err
}

func (d *DataInteractor) getDevicesData(token string, query *entities.Query) ([]entities.Data, error) {
	data := []entities.Data{}
	things, err := d.things.List(token)
	if err != nil {
		return data, err
	}

	if query.ThingID != "" {
		err := d.verifyAuthorization(token, query.ThingID)
		if err != nil {
			return nil, err
		}
	}

	data, err = d.getAllData(things, query)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (d *DataInteractor) getAllData(things []*btEntities.Thing, query *entities.Query) ([]entities.Data, error) {
	data := []entities.Data{}
	for _, t := range things {
		query.ThingID = t.ID
		fmt.Println(t.ID)
		thingData, err := d.DataStore.Get(query)
		if err != nil {
			return data, err
		}

		data = append(data, thingData...)
	}

	return data, nil
}

func (d *DataInteractor) verifyAuthorization(token, id string) error {
	things, err := d.things.List(token)
	if err != nil {
		return err
	}

	for _, t := range things {
		fmt.Println(t.ID)
		if t.ID == id {
			return nil
		}
	}

	return errors.New("user is not authorized to list thing's data")
}

func (d *DataInteractor) filterDataBySensorID(data []entities.Data, sensorID string) []entities.Data {
	if sensorID == "" {
		return data
	}

	id, err := strconv.Atoi(sensorID)
	if err != nil {
		d.logger.Errorf("Error when trying to parse ID from string to int")
		return data
	}

	filteredData := make([]entities.Data, 0)
	for _, v := range data {
		if v.Payload.SensorID == id {
			filteredData = append(filteredData, v)
		}
	}

	return filteredData
}
