package things

import (
	"errors"

	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
)

var (
	// ErrTokenEmpty is returned for empty access tokens.
	ErrTokenEmpty = errors.New("no access token provided")
	// ErrThingIDEmpty is returned for empty thing ID.
	ErrThingIDEmpty = errors.New("no thing ID provided")
)

// Lister specifies an API to use the service.
type Lister interface {
	List(token string) ([]*entities.Thing, error)
}

type service struct {
	proxy http.ThingProxy
}

// New creates a Things service instance.
func New(host string, port uint16, logger logging.Logger) Lister {
	return &service{
		proxy: http.NewThingProxy(logger, host, port),
	}
}

func (svc *service) List(token string) ([]*entities.Thing, error) {
	things, err := svc.proxy.List(token)
	if err != nil {
		return nil, err
	}
	return things, nil
}
