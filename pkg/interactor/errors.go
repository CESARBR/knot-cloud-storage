package interactor

import "errors"

var (
	// ErrUserNotAuthorized is returned when the user is not authorized to list the thing requested
	ErrUserNotAuthorized = errors.New("user is not authorized to list thing's data")

	// ErrDeviceIDNotProvided is returned when the user request to filter the data by sensorID, but don't provide the deviceID
	ErrDeviceIDNotProvided = errors.New("deviceID not provided")
)
