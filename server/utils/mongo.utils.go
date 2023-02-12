package utils

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	Config "github.com/oabraham1/open-blogger/server/config"
)

func Connect() (*mongo.Client, context.Context, error) {
	config := Config.NewConfig()
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DB_URL))
	if err != nil {
		return nil, nil, err
	}
	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, nil, err
	}
	log.Println("Connected to MongoDB!")

	return client, ctx, nil
}
