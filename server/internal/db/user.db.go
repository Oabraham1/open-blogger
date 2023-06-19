package db

import (
	"github.com/Oabraham1/open-blogger/server/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDB struct {
	MongoDB 	*MongoDB
	Collection 	*mongo.Collection
}

func NewUserDB() *UserDB {
	mongoDB := ConnectToMongo()
	defer mongoDB.Disconnect()
	return &UserDB{
		MongoDB: mongoDB,
	}
}

func (userDB *UserDB) CreateUser(user models.User) (primitive.ObjectID, error) {
	userID, err := userDB.MongoDB.InsertSingleDocument(userDB.Collection, user)
	if err != nil {
		return primitive.NewObjectID(), err
	}
	return userID.InsertedID.(primitive.ObjectID), nil
}

func (userDB *UserDB) FindUser(filter interface{}) (*models.User, error) {
	user := userDB.MongoDB.FindSingleDocument(userDB.Collection, filter)
	if user.Err() != nil {
		return nil, user.Err()
	}

	var userModel models.User
	err := user.Decode(&userModel)
	if err != nil {
		return nil, err
	}
	return &userModel, nil
}

func (userDB *UserDB) UpdateUser(filter interface{}, update interface{}) error {
	_, err := userDB.MongoDB.UpdateSingleDocument(userDB.Collection, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (userDB *UserDB) DeleteUser(filter interface{}) error {
	_, err := userDB.MongoDB.DeleteSingleDocument(userDB.Collection, filter)
	if err != nil {
		return err
	}
	return nil
}