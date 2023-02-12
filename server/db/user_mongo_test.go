package db

import (
	"context"
	"testing"

	user "github.com/oabraham1/open-blogger/server/internal/models/user"
	Utils "github.com/oabraham1/open-blogger/server/utils"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser() *user.User {
	return user.NewUser(
		"100",
		"testUser",
		"test",
		"user",
		"testuser@email.com",
		"http://testuser.com/image",
	)
}

func Setup() (*mongo.Collection, *user.User, error) {
	user := CreateUser()
	client, _, err := Utils.Connect()
	if err != nil {
		return nil, nil, err
	}
	collection := client.Database("test").Collection("users")
	if err != nil {
		return nil, nil, err
	}
	return collection, user, nil
}

func TearDown(collection *mongo.Collection, ctx context.Context) error {
	err := collection.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

func TestAddNewDBEntry(t *testing.T) {
	collection, user, err := Setup()
	require.NoError(t, err)

	_, err = AddNewDBEntry(collection, user, map[string]interface{}{"email": user.Email})
	require.NoError(t, err)

	_, err = AddNewDBEntry(collection, user, map[string]interface{}{"email": user.Email})
	require.Error(t, err)

	err = TearDown(collection, context.Background())
	require.NoError(t, err)
}

func TestFindUser(t *testing.T) {
	collection, user, err := Setup()
	require.NoError(t, err)

	_, err = AddNewDBEntry(collection, user, map[string]interface{}{"email": user.Email})
	require.NoError(t, err)
	_, err = FindDBEntry(collection, map[string]interface{}{"email": user.Email})
	require.NoError(t, err)

	_, err = FindDBEntry(collection, map[string]interface{}{"email": "notAnEmail"})
	require.Error(t, err)

	err = TearDown(collection, context.Background())
	require.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	collection, user, err := Setup()
	require.NoError(t, err)

	_, err = AddNewDBEntry(collection, user, map[string]interface{}{"email": user.Email})
	require.NoError(t, err)

	err = UpdateDBEntry(collection, map[string]interface{}{"email": user.Email}, bson.M{"$set": bson.M{"email": "newEmail"}})
	require.NoError(t, err)

	err = UpdateDBEntry(collection, map[string]interface{}{"email": "notAnEmail"}, bson.M{"$set": bson.M{"email": "newEmail"}})
	require.Error(t, err)

	err = TearDown(collection, context.Background())
	require.NoError(t, err)
}

func TestDelete(t *testing.T) {
	collection, user, err := Setup()
	require.NoError(t, err)

	_, err = AddNewDBEntry(collection, user, map[string]interface{}{"email": user.Email})
	require.NoError(t, err)

	err = DeleteDBEntry(collection, map[string]interface{}{"email": user.Email})
	require.NoError(t, err)

	err = DeleteDBEntry(collection, map[string]interface{}{"email": "notAnEmail"})
	require.Error(t, err)

	err = TearDown(collection, context.Background())
	require.NoError(t, err)
}
