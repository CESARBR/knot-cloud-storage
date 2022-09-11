package interactor

import (
	"fmt"

	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
)

// Save inserts data to the storage, if it doesn't exist already
func (d *DataInteractor) Save(token, id string, data []entities.Payload) error {
	_, ok := d.tokenCache[token]
	if ok {
		return nil
	} else {
		_, err := d.things.List(token)
		if err != nil {
			return fmt.Errorf("%v: %w", ErrValidToken, err)
		}
		// Updates the cache to signal the existence of the token.
		d.tokenCache[token] = struct{}{}
		for _, dt := range data {
			err = d.DataStore.Save(entities.Data{From: id, Payload: dt})
			if err != nil {
				return fmt.Errorf("error saving data %v: %w", data, err)
			}
		}

		return nil
	}
}
