package interactor

func (d *DataInteractor) Delete(token, deviceID string) error {
	_, err := d.things.List(token)
	if err != nil {
		return err
	}

	err = d.DataStore.Delete(deviceID)
	if err != nil {
		d.logger.Error(err)
		return err
	}
	return nil
}
