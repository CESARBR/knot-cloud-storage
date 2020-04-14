package entities

import (
	"time"
)

type Data struct {
	From      string    `bson:"from" json:"from"`
	Payload   Payload   `bson:"payload" json:"payload"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type Payload struct {
	SensorId int         `bson:"sensorId" json:"sensorId"`
	Value    interface{} `bson:"value" json:"value"`
}
