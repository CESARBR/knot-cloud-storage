package entities

import (
	"time"
)

// Query represents the domain data query properties
type Query struct {
	ThingID    string
	SensorID   string
	Order      int
	Skip       int
	Take       int
	StartDate  time.Time
	FinishDate time.Time
}
