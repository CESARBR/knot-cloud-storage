package interactor

import "errors"

var (
	// ErrTokenEmpty is returned when no authorization token is provided
	ErrTokenEmpty = errors.New("authorization token not provided")

	// ErrUserNotAuthorized is returned when the user is not authorized to list the thing requested
	ErrUserNotAuthorized = errors.New("user is not authorized to list thing's data")

	// ErrDeviceIDNotProvided is returned when the user request to filter the data by sensorID, but don't provide the deviceID
	ErrDeviceIDNotProvided = errors.New("deviceID not provided")

	// ErrValidToken is returned when the service to valid the token returns an error
	ErrValidToken = errors.New("failed to valid token")

	// ErrToParseString is returned when the parse from string to int of the sensorID fails
	ErrToParseString = errors.New("failed to parse ID from string to int")
)
