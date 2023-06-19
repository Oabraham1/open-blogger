package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client 	*mongo.Client
	error  	error
	ctx   	context.Context
}

var mongoUri = os.Getenv("MONGO_URI")

func ConnectToMongo() *MongoDB {
	mongoDB := &MongoDB{}
	mongoDB.ctx = context.Background()
	mongoDB.client, mongoDB.error = mongo.Connect(mongoDB.ctx, options.Client().ApplyURI(mongoUri))

	return mongoDB
}

func (mongoDB *MongoDB) Disconnect() {
	mongoDB.client.Disconnect(mongoDB.ctx)
}

func (mongoDB *MongoDB) InsertSingleDocument(collection  *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error){
	return collection.InsertOne(mongoDB.ctx, document)
}

func (mongoDB *MongoDB) FindSingleDocument(collection  *mongo.Collection, filter interface{}) *mongo.SingleResult {
	return collection.FindOne(mongoDB.ctx, filter)
}

func (mongoDB *MongoDB) FindMultipleDocuments(collection  *mongo.Collection, filter interface{}) (*mongo.Cursor, error) {
	return collection.Find(mongoDB.ctx, filter)
}

func (mongoDB *MongoDB) UpdateSingleDocument(collection  *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return collection.UpdateOne(mongoDB.ctx, filter, update)
}

func (mongoDB *MongoDB) DeleteSingleDocument(collection  *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	return collection.DeleteOne(mongoDB.ctx, filter)
}

func (mongoDB *MongoDB) DeleteMultipleDocuments(collection  *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	return collection.DeleteMany(mongoDB.ctx, filter)
}