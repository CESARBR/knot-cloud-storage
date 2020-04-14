package users

import (
	"fmt"

	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
)

// Authenticator specifies an API to use the service.
type Authenticator interface {
	Authenticate(token string) error
}

type service struct {
	url    string
	logger logging.Logger
}

// New creates a Users service instance.
func New(host string, port int, logger logging.Logger) Authenticator {
	return &service{
		url:    fmt.Sprintf("http://%s:%d/%s", host, port, "users"),
		logger: logger,
	}
}

// Authenticate sends a request to the cloud to verify the access token.
func (svc *service) Authenticate(token string) error {
	err := svc.sendAuthRequest(token)
	if err != nil {
		svc.logger.Errorf("authentication request failed: %v", err)
		return err
	}
	return nil
}
