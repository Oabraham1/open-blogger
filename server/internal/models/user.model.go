package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID 				primitive.ObjectID `json:"id" bson:"_id"`
	Username 		string `json:"username" bson:"username"`
	Password 		string `json:"password" bson:"password"`
	Email 			string `json:"email" bson:"email"`
	FirstName 		string `json:"first_name" bson:"first_name"`
	LastName 		string `json:"last_name" bson:"last_name"`
}