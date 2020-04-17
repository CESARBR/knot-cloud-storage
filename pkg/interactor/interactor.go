package interactor

import (
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/CESARBR/knot-cloud-storage/pkg/things"
)

// Interactor is an interface that defines the data's use cases operations
type Interactor interface {
	Save(token, id string, data []entities.Payload, time time.Time) error
	List(token string, query *entities.Query) ([]entities.Data, error)
	Delete(token, deviceID string) error
}

// DataInteractor represents the data layer interactor structure
type DataInteractor struct {
	things    things.Lister
	DataStore *data.Store
	logger    logging.Logger
}

// NewDataInteractor creates a new data interactor instance
func NewDataInteractor(things things.Lister, dataStore *data.Store, logger logging.Logger) Interactor {
	return &DataInteractor{things, dataStore, logger}
}
