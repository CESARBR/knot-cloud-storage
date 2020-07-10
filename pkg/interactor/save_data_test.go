package interactor

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	// ErrSaveDataDB is returned when the database fails to save data
	errSaveData = errors.New("failed to save data")
)

type SaveDataTestCase struct {
	name             string
	authorization    string
	deviceID         string
	payload          []entities.Payload
	data             entities.Data
	fakeLogger       *mocks.FakeLogger
	fakeDataStore    *mocks.FakeDataStore
	fakeThingService *mocks.FakeThingService
	expectedError    error
}

var sdCases = []SaveDataTestCase{
	{
		"authorization token not provided",
		"",
		"",
		nil,
		entities.Data{},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{},
		&mocks.FakeThingService{ReturnErr: ErrTokenEmpty},
		ErrTokenEmpty,
	},
	{
		"data successfully saved on the database",
		"authorization-token",
		"fc3fcf912d0c290a",
		[]entities.Payload{{SensorID: 1, Value: 999}},
		entities.Data{From: "fc3fcf912d0c290a", Payload: entities.Payload{SensorID: 1, Value: 999}},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{},
		&mocks.FakeThingService{},
		nil,
	},
	{
		"failed to save the data on the database",
		"authorization-token",
		"fc3fcf912d0c290a",
		[]entities.Payload{{SensorID: 1, Value: 999}},
		entities.Data{From: "fc3fcf912d0c290a", Payload: entities.Payload{SensorID: 1, Value: 999}},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{
			ReturnErr: errSaveData,
		},
		&mocks.FakeThingService{},
		errSaveData,
	},
}

func TestSaveData(t *testing.T) {
	for _, tc := range sdCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeDataStore.
				On("Save", tc.data).
				Return(tc.fakeDataStore.ReturnErr).
				Maybe()
			tc.fakeThingService.
				On("List", tc.authorization).
				Return(tc.fakeThingService.Things, tc.fakeThingService.ReturnErr).
				Maybe()
		})

		dataInteractor := NewDataInteractor(tc.fakeThingService, tc.fakeDataStore, tc.fakeLogger)
		err := dataInteractor.Save(tc.authorization, tc.deviceID, tc.payload)

		assert.EqualValues(t, errors.Is(err, tc.expectedError), true)

		tc.fakeDataStore.AssertExpectations(t)
		tc.fakeThingService.AssertExpectations(t)
	}
}
