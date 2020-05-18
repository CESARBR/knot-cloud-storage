package interactor

import (
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

	data, err = d.filterDataBySensorID(data, query.ThingID, query.SensorID)
	if err != nil {
		return []entities.Data{}, err
	}

	return data, err
}

func (d *DataInteractor) getDevicesData(token string, query *entities.Query) ([]entities.Data, error) {
	data := []entities.Data{}
	things, err := d.things.List(token)
	if err != nil {
		return data, fmt.Errorf("%s: %v", ErrValidToken, err)
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
			return data, fmt.Errorf("error getting data: %w", err)
		}

		data = append(data, thingData...)
	}

	return data, nil
}

func (d *DataInteractor) verifyAuthorization(token, id string) error {
	things, err := d.things.List(token)
	if err != nil {
		return fmt.Errorf("%s: %v", ErrValidToken, err)
	}

	for _, t := range things {
		fmt.Println(t.ID)
		if t.ID == id {
			return nil
		}
	}

	return ErrUserNotAuthorized
}

func (d *DataInteractor) filterDataBySensorID(data []entities.Data, thingID string, sensorID string) ([]entities.Data, error) {
	if sensorID == "" {
		return data, nil
	}

	if thingID == "" {
		return nil, ErrDeviceIDNotProvided
	}

	id, err := strconv.Atoi(sensorID)
	if err != nil {
		d.logger.Errorf("failed to parse ID from string to int")
		return nil, ErrToParseString
	}

	filteredData := make([]entities.Data, 0)
	for _, v := range data {
		if v.Payload.SensorID == id {
			filteredData = append(filteredData, v)
		}
	}

	return filteredData, nil
}
