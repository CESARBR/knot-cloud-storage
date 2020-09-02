package main

import (
	"bytes"
	"net/http"
	"github.com/google/go-querystring/query"
	"time"
)

func sendRequest(svc *requestInfo) (*http.Response, error) {
	values, err := query.Values(svc.options)
	failOnError(err, "Failed to get query values")

	queryString := "?" + values.Encode()
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(svc.method, svc.url+queryString, bytes.NewBuffer(svc.data))
	failOnError(err, "Failed to do a HTTP request")

	req.Header.Set("auth_token", svc.authorization)
	req.Header.Set("Content-Type", svc.contentType)

	return client.Do(req)
}
