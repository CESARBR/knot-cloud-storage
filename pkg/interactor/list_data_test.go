package interactor

import (
	"errors"
	"testing"
	"time"

	thingentities "github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	// ErrDeleteDataDB is returned when the database fails to get data
	errGetData = errors.New("failed to get data")
)

type ListDataTestCase struct {
	name             string
	authorization    string
	query            *entities.Query
	fakeLogger       *mocks.FakeLogger
	fakeDataStore    *mocks.FakeDataStore
	fakeThingService *mocks.FakeThingService
	expectedError    error
}

var ldCases = []ListDataTestCase{
	{
		"authorization token not provided",
		"",
		&entities.Query{},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{},
		&mocks.FakeThingService{
			ReturnErr: errTokenEmpty,
		},
		errTokenEmpty,
	},
	{
		"successful get on the database",
		"authorization-token",
		&entities.Query{},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{
			Data: []entities.Data{
				{
					From: "fc3fcf912d0c290a",
					Payload: entities.Payload{
						SensorID: 1,
						Value:    3,
					},
					Timestamp: time.Now(),
				},
			},
		},
		&mocks.FakeThingService{
			Things: []*thingentities.Thing{
				{
					ID:    "fc3fcf912d0c290a",
					Token: "token",
					Name:  "thing",
					Schema: []thingentities.Schema{
						{
							SensorID:  1,
							ValueType: 3,
							Unit:      0,
							TypeID:    65521,
							Name:      "Test",
						},
					},
				},
			},
		},
		nil,
	},
	{
		"successful get on the database by thingID",
		"authorization-token",
		&entities.Query{
			ThingID: "fc3fcf912d0c290a",
		},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{
			Data: []entities.Data{
				{
					From: "fc3fcf912d0c290a",
					Payload: entities.Payload{
						SensorID: 1,
						Value:    3,
					},
					Timestamp: time.Now(),
				},
			},
		},
		&mocks.FakeThingService{
			Things: []*thingentities.Thing{
				{
					ID:    "fc3fcf912d0c290a",
					Token: "token",
					Name:  "thing",
					Schema: []thingentities.Schema{
						{
							SensorID:  1,
							ValueType: 3,
							Unit:      0,
							TypeID:    65521,
							Name:      "Test",
						},
					},
				},
			},
		},
		nil,
	},
	{
		"successful get on the database filtered by sensorID",
		"authorization-token",
		&entities.Query{
			ThingID:  "fc3fcf912d0c290a",
			SensorID: "1",
		},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{
			Data: []entities.Data{
				{
					From: "fc3fcf912d0c290a",
					Payload: entities.Payload{
						SensorID: 1,
						Value:    3,
					},
					Timestamp: time.Now(),
				},
			},
		},
		&mocks.FakeThingService{
			Things: []*thingentities.Thing{
				{
					ID:    "fc3fcf912d0c290a",
					Token: "token",
					Name:  "thing",
					Schema: []thingentities.Schema{
						{
							SensorID:  1,
							ValueType: 3,
							Unit:      0,
							TypeID:    65521,
							Name:      "Test",
						},
					},
				},
			},
		},
		nil,
	},
	{
		"deviceID not provided to filter by sensorID",
		"authorization-token",
		&entities.Query{
			SensorID: "1",
		},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{
			Data: []entities.Data{
				{
					From: "fc3fcf912d0c290a",
					Payload: entities.Payload{
						SensorID: 1,
						Value:    3,
					},
					Timestamp: time.Now(),
				},
			},
		},
		&mocks.FakeThingService{},
		ErrDeviceIDNotProvided,
	},
	{
		"user not authorized to list thing's data",
		"authorization-token",
		&entities.Query{
			ThingID: "fc4fcf525d0c693b",
		},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{
			Data: []entities.Data{
				{
					From: "fc4fcf525d0c693b",
					Payload: entities.Payload{
						SensorID: 1,
						Value:    3,
					},
					Timestamp: time.Now(),
				},
			},
		},
		&mocks.FakeThingService{
			Things: []*thingentities.Thing{
				{
					ID:    "fc3fcf912d0c290a",
					Token: "token",
					Name:  "thing",
					Schema: []thingentities.Schema{
						{
							SensorID:  1,
							ValueType: 3,
							Unit:      0,
							TypeID:    65521,
							Name:      "Test",
						},
					},
				},
			},
		},
		ErrUserNotAuthorized,
	},
	{
		"failed to get the data on the database",
		"authorization-token",
		&entities.Query{},
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{
			Data:      []entities.Data{},
			ReturnErr: errGetData,
		},
		&mocks.FakeThingService{
			Things: []*thingentities.Thing{
				{
					ID:    "fc3fcf912d0c290a",
					Token: "token",
					Name:  "thing",
					Schema: []thingentities.Schema{
						{
							SensorID:  1,
							ValueType: 3,
							Unit:      0,
							TypeID:    65521,
							Name:      "Test",
						},
					},
				},
			},
		},
		errGetData,
	},
}

func TestListData(t *testing.T) {
	for _, tc := range ldCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeDataStore.
				On("Get", tc.query).
				Return(tc.fakeDataStore.Data, tc.fakeDataStore.ReturnErr).
				Maybe()
			tc.fakeThingService.
				On("List", tc.authorization).
				Return(tc.fakeThingService.Things, tc.fakeThingService.ReturnErr).
				Maybe()
		})

		dataInteractor := NewDataInteractor(tc.fakeThingService, tc.fakeDataStore, tc.fakeLogger)
		_, err := dataInteractor.List(tc.authorization, tc.query)

		assert.EqualValues(t, errors.Is(err, tc.expectedError), true)

		tc.fakeDataStore.AssertExpectations(t)
		tc.fakeThingService.AssertExpectations(t)
	}
}
