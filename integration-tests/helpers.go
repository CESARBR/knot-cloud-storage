package main

import (
	"encoding/json"
	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/internal/config"
	"log"
	"regexp"
	"strings"
	"time"
)

type connection struct {
	Server, Port string
}

type thing struct {
	ID string
	Data []entities.Payload
}

type userCredential struct {
	Email string
	Password string
	Server string
	Port string
}

type userToken struct {
	Token string
}

func evalOsArgs(osInput []string) *userCredential {
	return &userCredential{
		Email: getEmail(osInput[1]),
		Password: getPassword(osInput[2]),
		Server: getServer(osInput[3]),
		Port: getPort(osInput[4])}
}

func getEmail(email string) string {
	regex := regexp.MustCompile(`[a-z0-9]*\@[a-z]*\.[a-z]*`)
	eval := regex.MatchString(email)
	if eval != true {
		log.Println("Email invalid. Switching to default email...")
		return "test@test.com"
	}
	return email
}

func getDates(initialDate string, finalDate time.Time) (time.Time, time.Time) {
	startDate, _ := time.Parse("2006-01-02 15:04:05", initialDate)
	finishDate := finalDate
	return startDate, finishDate
}

func getLogger(value string) *config.Logger {
	return &config.Logger{Level: value}
}

func getNullPayloadData(payload *thing) {
	if len(payload.Data) == 0 {
		payload.Data = []entities.Payload{{SensorID: 0, Value: nil},}
	}
}

func getPassword(password string) string {
	eval := strings.EqualFold(password, "")
	if eval != false {
		log.Println("Password invalid. Switching to default password...")
		return "abcdef"
	}
	return password
}

func getPayload(id string, data []entities.Payload) (*thing, []byte) {
	payload := &thing{ID: id, Data: data}
	body, err := json.Marshal(payload)
	failOnError(err, "Failed to convert to JSON")
	return payload, body
}

func getPort(port string) string {
	eval := strings.EqualFold(port, "")
	if eval != false {
		log.Println("Port invalid. Switching to default port...")
		return "8180"
	}
	return port
}

func getQuery(thingID string, sensorID string, order int, skip int64, take int64, startDate time.Time, finishDate time.Time) *entities.Query {
	return &entities.Query{
		ThingID: thingID,
		SensorID: sensorID,
		Order: order,
		Skip: skip,
		Take: take,
		StartDate: startDate,
		FinishDate: finishDate,
	}
}

func getServer(server string) string {
	eval := strings.EqualFold(server, "")
	if eval != false {
		log.Println("Server invalid. Switching to default server...")
		return "localhost"
	}
	return server
}

func getURL(usrCred *userCredential, folder string) string {
	str := []string{"http://", usrCred.Server, ":", usrCred.Port, "/", folder}
	return strings.Join(str, "")
}

func getUserCredentials(osInput []string) *userCredential {
	if len(osInput) > 1 {
		return evalOsArgs(osInput)
	}
	return &userCredential{
		Email: "test@test.com",
		Password: "abcdef",
		Server: "localhost",
		Port: "8180"}
}

func getUserJSON(usrCred *userCredential) []byte {
	data := map[string]string{"email": usrCred.Email, "password": usrCred.Password}
	userJSON, err := json.Marshal(data)
	failOnError(err, "Failed to make user JSON")
	return userJSON
}
