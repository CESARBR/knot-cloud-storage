package interactor

import (
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/CESARBR/knot-cloud-storage/pkg/users"
)

// Interactor is an interface that defines the data's use cases operations
type Interactor interface {
	Save(token, id string, data []entities.Payload, time time.Time) error
	List(token string, query *entities.Query) ([]entities.Data, error)
}

// DataInteractor represents the data layer interactor structure
type DataInteractor struct {
	UsersService users.Authenticator
	DataStore    *data.DataStore
	logger       logging.Logger
}

// NewDataInteractor creates a new data interactor instance
func NewDataInteractor(users users.Authenticator, dataStore *data.DataStore, logger logging.Logger) *DataInteractor {
	return &DataInteractor{users, dataStore, logger}
}

// Authenticate verifies if the access token is valid.
func (d *DataInteractor) Authenticate(token string) error {
	err := d.UsersService.Authenticate(token)
	if err != nil {
		return err
	}
	return nil
}