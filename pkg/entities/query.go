package entities

import (
	"time"
)

// Query represents the domain data query properties
type Query struct {
	ThingID    string
	SensorID   string
	Order      int
	Skip       int64
	Take       int64
	StartDate  time.Time
	FinishDate time.Time
}
