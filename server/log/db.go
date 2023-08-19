package logger

import (
	"context"

	"github.com/Oabraham1/open-blogger/server/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	config, err := util.LoadConfig("../../..")
	if err != nil {
		return nil, err
	}
	mongoURI := config.MongoURI

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
