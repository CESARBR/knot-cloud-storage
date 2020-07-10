package interactor

import (
	"fmt"
	"strconv"

	btEntities "github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
)

// List provides all the thing's data owned by an user
func (d *DataInteractor) List(token string, query *entities.Query) ([]entities.Data, error) {
	if token == "" {
		return nil, ErrTokenEmpty
	}

	things, err := d.things.List(token)
	if err != nil {
		return nil, fmt.Errorf("error getting list of things: %w", err)
	}

	if query.ThingID != "" {
		things, err = d.verifyThingID(things, query.ThingID)
		if err != nil {
			return nil, fmt.Errorf("thing ID not in user things list: %w", err)
		}
	}

	data, err := d.getAllData(things, query)
	if err != nil {
		return nil, fmt.Errorf("error getting the data: %w", err)
	}

	if query.SensorID != "" {
		data, err = d.filterDataBySensorID(data, query.SensorID)
		if err != nil {
			return nil, fmt.Errorf("error filtering data by sensor ID: %w", err)
		}
	}

	return data, nil
}

func (d *DataInteractor) verifyThingID(things []*btEntities.Thing, thingID string) ([]*btEntities.Thing, error) {
	for _, t := range things {
		if t.ID == thingID {
			return []*btEntities.Thing{{ID: thingID}}, nil
		}
	}
	return nil, ErrUserNotAuthorized
}

func (d *DataInteractor) getAllData(things []*btEntities.Thing, query *entities.Query) ([]entities.Data, error) {
	data := []entities.Data{}
	for _, t := range things {
		query.ThingID = t.ID
		thingData, err := d.DataStore.Get(query)
		if err != nil {
			return nil, fmt.Errorf("error getting data: %w", err)
		}
		data = append(data, thingData...)
	}
	return data, nil
}

func (d *DataInteractor) filterDataBySensorID(data []entities.Data, sensorID string) ([]entities.Data, error) {
	id, err := strconv.Atoi(sensorID)
	if err != nil {
		return nil, fmt.Errorf("error parsing ID from string to int: %w", err)
	}

	sensorData := []entities.Data{}
	for _, dt := range data {
		if dt.Payload.SensorID == id {
			sensorData = append(sensorData, dt)
		}
	}

	return sensorData, nil
}
