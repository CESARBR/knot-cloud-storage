package interactor

import (
	"errors"
	"testing"

	thingentities "github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	// ErrDeleteDataDB is returned when the database fails to delete data
	errDeleteData = errors.New("failed to delete data")
)

type DeleteDataTestCase struct {
	name             string
	authorization    string
	deviceID         string
	fakeLogger       *mocks.FakeLogger
	fakeDataStore    *mocks.FakeDataStore
	fakeThingService *mocks.FakeThingService
	expectedError    error
}

var ddCases = []DeleteDataTestCase{
	{
		"authorization token not provided",
		"",
		"",
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{},
		&mocks.FakeThingService{
			ReturnErr: errTokenEmpty,
		},
		errTokenEmpty,
	},
	{
		"successful delete on the database",
		"authorization-token",
		"fc3fcf912d0c290a",
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{},
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
		"failed to delete on the database",
		"authorization-token",
		"fc3fcf912d0c290a",
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{
			ReturnErr: errDeleteData,
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
		errDeleteData,
	},
	{
		"successful delete on the database without deviceID",
		"authorization-token",
		"",
		&mocks.FakeLogger{},
		&mocks.FakeDataStore{},
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
}

func TestDeleteData(t *testing.T) {
	for _, tc := range ddCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeDataStore.
				On("Delete", tc.deviceID).
				Return(tc.fakeDataStore.ReturnErr).
				Maybe()
			tc.fakeThingService.
				On("List", tc.authorization).
				Return(tc.fakeThingService.Things, tc.fakeThingService.ReturnErr).
				Maybe()
		})

		dataInteractor := NewDataInteractor(tc.fakeThingService, tc.fakeDataStore, tc.fakeLogger)
		err := dataInteractor.Delete(tc.authorization, tc.deviceID)

		assert.EqualValues(t, errors.Is(err, tc.expectedError), true)

		tc.fakeDataStore.AssertExpectations(t)
		tc.fakeThingService.AssertExpectations(t)
	}
}
