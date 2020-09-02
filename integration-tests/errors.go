package main

import (
	"errors"
	"log"
	"net/http"
)

var (
	errEmptyToken = errors.New("It should have a non-empty string for user token")
	errInvalidToken = errors.New("should fail if provided token is invalid")
	errQueryResult = errors.New("should fail if the last query output is different from the query input")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func mapErrorFromStatusCode(code int) error {
	if code != http.StatusCreated {
		switch code {
		case http.StatusUnauthorized:
			return errInvalidToken
		}
	}
	return nil
}
