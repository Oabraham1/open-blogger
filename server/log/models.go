package logger

import "time"

type ErrorLogger struct {
	Body           string    `json:"body" bson:"body"`
	LoggedAt       time.Time `json:"loggedAt" bson:"loggedAt"`
	PointOfFailure string    `json:"pointOfFailure" bson:"pointOfFailure"`
}
