package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/Oabraham1/open-blogger/server/util"
)

func LogError(body interface{}, pointOfFailure string) {
	errorLog := ErrorLogger{
		Body:           fmt.Sprintf("%v", body),
		LoggedAt:       time.Now(),
		PointOfFailure: pointOfFailure,
	}

	SendErrorLogToDatabase(errorLog)
}

func SendErrorLogToDatabase(errorObject ErrorLogger) {
	config, err := util.LoadConfig("../../..")
	if err != nil {
		return
	}

	mongoClient, err := ConnectToMongoDB()
	if err != nil {
		return
	}

	defer mongoClient.Disconnect(context.TODO())

	environment := config.Environment

	collection := mongoClient.Database("openBloggerDB").Collection(environment)
	_, err = collection.InsertOne(context.Background(), errorObject)
	if err != nil {
		return
	}
}
