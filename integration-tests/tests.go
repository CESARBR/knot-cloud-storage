package main

import (
	"fmt"
	"log"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/google/go-cmp/cmp"
)

type requestInfo struct {
	method string
	url string
	authorization string
	contentType string
	data []byte
	options *requestOptions
}

type requestOptions struct {
	limit int
	offset int
}

func setRequest(usrToken string, body[]byte) *requestInfo {
	return &requestInfo{
		method: "POST",
		url: "http://localhost:80/data",
		authorization: usrToken,
		contentType: "application/json",
		data: body,
		options: &requestOptions{limit: 100, offset: 0},
	}
}

func showQueryResult(output []entities.Data) {
	fmt.Printf("\n\n")
	fmt.Println("THING_ID SENSOR_ID VALUE TIMESTAMP")
	for i := range output {
		fmt.Println(output[i])
	}
}

func startTest(msg error, testFunction error) {
	status := "PASS"
	err := testFunction
	if err != nil {
		status = "FAIL"
	}
	log.Printf("%s: %s", status, msg)

}

// If the test fails, que result is not shown on the screen
func testQueryOutput(output []entities.Data, payload *thing) error {
	err := errQueryResult
	equalThingID := cmp.Equal(output[len(output) - 1].From, payload.ID)
	equalPayload := cmp.Equal(output[len(output) - 1].Payload, payload.Data[len(payload.Data) - 1])
	if equalThingID == equalPayload {
		err = nil
	}
	return err
}

// When the test fails, the payload is not published in the MongoDB
func testUserToken(usrToken string, body []byte) error {
	if usrToken == "" {
		return errEmptyToken
	}

	svc := setRequest(usrToken, body)
	rqs, _ := sendRequest(svc)
	defer rqs.Body.Close()

	err := mapErrorFromStatusCode(rqs.StatusCode)
	if err != nil {
		return err
	}

	publishData(usrToken, body)
	return nil
}

