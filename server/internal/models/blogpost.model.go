package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BlogPost struct {
	ID 			primitive.ObjectID 	`json:"id" bson:"_id"`
	Title 		string 				`json:"title" bson:"title"`
	Content 	string 				`json:"content" bson:"content"`
	Author 		string 				`json:"author" bson:"author"`
	AuthorID 	string 				`json:"author_id" bson:"author_id"`
}