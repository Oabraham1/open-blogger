package db

import (
	"testing"

	"github.com/Oabraham1/open-blogger/server/internal/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("CreateUser", func(mt *mtest.T) {
		userCollection := mt.Coll
		userDB := &UserDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: userCollection,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		userID, err := userDB.CreateUser(models.User{
			Username:  "Test Username",
			Email:     "Test Email",
			FirstName: "Test First Name",
			LastName:  "Test Last Name",
		})
		require.NoError(t, err)
		require.NotNil(t, userID)

		mt.ClearMockResponses()

	})
}

func TestFindUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("FindUser", func(mt *mtest.T) {
		userCollection := mt.Coll
		userDB := &UserDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: userCollection,
		}
		expectedUser := models.User{
			Username:  "Test Username",
			Email:     "Test Email",
			FirstName: "Test First Name",
			LastName:  "Test Last Name",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedUser.ID},
			{Key: "username", Value: expectedUser.Username},
			{Key: "email", Value: expectedUser.Email},
			{Key: "first_name", Value: expectedUser.FirstName},
			{Key: "last_name", Value: expectedUser.LastName},
		}))

		user, err := userDB.FindUser(bson.M{"ID": expectedUser.Username})
		require.NoError(t, err)
		require.NotNil(t, user)

		mt.ClearMockResponses()

	})
}

func TestUpdateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("UpdateUser", func(mt *mtest.T) {
		userCollection := mt.Coll
		userDB := &UserDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: userCollection,
		}
		expectedUser := models.User{
			Username:  "Test Username",
			Email:     "Test Email",
			FirstName: "Test First Name",
			LastName:  "Test Last Name",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := userDB.UpdateUser(bson.M{"ID": expectedUser.Username}, bson.M{"$set": bson.M{"username": expectedUser.Username}})
		require.NoError(t, err)

		mt.ClearMockResponses()

	})
}

func TestDeleteUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("DeleteUser", func(mt *mtest.T) {
		userCollection := mt.Coll
		userDB := &UserDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: userCollection,
		}
		expectedUser := models.User{
			Username:  "Test Username",
			Email:     "Test Email",
			FirstName: "Test First Name",
			LastName:  "Test Last Name",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := userDB.DeleteUser(bson.M{"ID": expectedUser.Username})
		require.NoError(t, err)

		mt.ClearMockResponses()

	})
}
