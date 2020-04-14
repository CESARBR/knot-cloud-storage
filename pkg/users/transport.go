package users

import (
	"errors"
	"net/http"
	"time"
)

var (
	// ErrBadRequest is returned when query parameters are missing or invalid.
	ErrBadRequest = errors.New("failed due to malformed query parameters")

	// ErrTokenInvalid is returned when the access token isn't registered in the cloud.
	ErrTokenInvalid = errors.New("invalid access token provided")

	// ErrServerUnexpected is returned when the server is responsible for the error.
	ErrServerUnexpected = errors.New("unexpected server-side error occurred")

	// ErrUnknown is returned when the error is none of the above.
	ErrUnknown = errors.New("unknown error")
)

func (svc *service) sendAuthRequest(token string) error {
	req, err := http.NewRequest("GET", svc.url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return svc.getResponse(resp.StatusCode)
}

func (svc *service) getResponse(statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusForbidden:
		return ErrTokenInvalid
	case http.StatusInternalServerError:
		return ErrServerUnexpected
	default:
		return ErrUnknown
	}
}
