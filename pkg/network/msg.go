package network

import (
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
)

// DataPublish represents the incoming publish data command
type DataPublish struct {
	ID   string             `json:"id"`
	Data []entities.Payload `json:"data"`
}
