package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/mock"
)

// FakeThingService represents a mocking type for the thing's service
type FakeThingService struct {
	mock.Mock
	ReturnErr error
	Things    []*entities.Thing
}

// List provides a mock function from the thing's service
func (fts *FakeThingService) List(token string) ([]*entities.Thing, error) {
	ret := fts.Called(token)
	return ret.Get(0).([]*entities.Thing), ret.Error(1)
}
