package main

import (
	"bytes"
	"encoding/json"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"net/http"
	"os"
	"time"
)

func main() {

	usr := getUserCredentials(os.Args)
	buff := getUserJSON(usr)
	url := getURL(usr, "users")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buff))
	failOnError(err, "Failed to do a HTTP POST")
	resp.Body.Close()

	// Make a token creation POST request
	tokenURL := getURL(usr, "tokens")
	tokenResp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(buff))
	failOnError(err, "Failed to do a HTTP POST")

	token := &userToken{}
	json.NewDecoder(tokenResp.Body).Decode(token)
	tokenResp.Body.Close()

	config := getLogger("debug")
	logrus := logging.NewLogrus(config.Level, false)
	logger := logrus.Get("Main")
	logger.Info("Starting storage test")

	payload, body := getPayload("fbe64efa6c7f717e", []entities.Payload{})
	startTest(errInvalidToken, testUserToken(token.Token, body))

	startDate, finishDate := getDates("2020-08-28 21:28:07", time.Now())
	getNullPayloadData(payload)

	newDatabase, ctx, client := setupDatabase("mongodb://localhost:27017/", "things_db", logger)
	defer disconnectClient(ctx, client)
	newStore := setupStore(newDatabase, logger, 0)

	query := getQuery(payload.ID, string(payload.Data[0].SensorID), 1, 0, 10, startDate, finishDate)
	output, err := newStore.Get(query)
	failOnError(err, "Failed to get data from MongoDB")
	startTest(errQueryResult, testQueryOutput(output, payload))
	showQueryResult(output)
}
