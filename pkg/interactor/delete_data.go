package interactor

func (d *DataInteractor) Delete(deviceID string) error {
	err := d.DataStore.Delete(deviceID)
	if err != nil {
		d.logger.Error(err)
		return err
	}
	return nil
}
