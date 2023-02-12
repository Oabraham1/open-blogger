package db

import (
	"errors"
	"log"

	Utils "github.com/oabraham1/open-blogger/server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddNewDBEntry(collection *mongo.Collection, input interface{}, filter interface{}) (interface{}, error) {
	_, ctx, err := Utils.Connect()
	if err != nil {
		log.Println("ERROR: ", err)
		return nil, err
	}

	if filter != nil {
		var result bson.M
		err = collection.FindOne(ctx, filter).Decode(&result)

		if err != nil && err.Error() != "mongo: no documents in result" {
			log.Println("ERROR: ", err)
			return nil, err
		}

		log.Println("RESULT: ", result)
		if len(result) > 0 {
			log.Println("User already exists in DB")
			return nil, errors.New("failed to insert user into DB")
		}
	}

	log.Println("Here.......")
	_, err = collection.InsertOne(ctx, input)
	if err != nil {
		log.Println("ERROR: ", err)
		return nil, err
	}

	return input, nil
}

func FindDBEntry(collection *mongo.Collection, filter interface{}) (interface{}, error) {
	_, ctx, err := Utils.Connect()
	if err != nil {
		log.Println("ERROR: ", err)
		return nil, err
	}

	var result interface{}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Println("ERROR: ", err)
		return nil, err
	}

	return result, nil
}

func UpdateDBEntry(collection *mongo.Collection, filter bson.M, update bson.M) error {
	_, ctx, err := Utils.Connect()
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	if result.MatchedCount == 0 {
		log.Println("ERROR: ", err)
		return errors.New("no user found")
	}
	return nil
}

func DeleteDBEntry(collection *mongo.Collection, filter bson.M) error {
	_, ctx, err := Utils.Connect()
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	err = collection.FindOneAndDelete(ctx, filter).Err()
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}
	return nil
}
