package mocks

import (
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/stretchr/testify/mock"
)

// FakeDataStore represents a mocking type for the data store
type FakeDataStore struct {
	mock.Mock
	ReturnErr error
	Data      []entities.Data
}

// Get provides a mock function to get data from the database
func (fds *FakeDataStore) Get(query *entities.Query) ([]entities.Data, error) {
	ret := fds.Called(query)
	return ret.Get(0).([]entities.Data), ret.Error(1)
}

// Save provides a mock function to save data on the database
func (fds *FakeDataStore) Save(data entities.Data) error {
	ret := fds.Called(data)
	return ret.Error(0)
}

// Delete provides a mock function to delete data from the database
func (fds *FakeDataStore) Delete(deviceID string) error {
	ret := fds.Called(deviceID)
	return ret.Error(0)
}
