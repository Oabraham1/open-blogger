package db

import (
	"log"
	"testing"

	"github.com/Oabraham1/open-blogger/server/internal/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateNewBlogBost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("CreateNewBlogPost", func(mt *mtest.T) {
		blogPostCollection := mt.Coll
		blogPostDB := &BlogPostDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: blogPostCollection,
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		blogPostID, err := blogPostDB.CreateNewBlogPost(models.BlogPost{
			ID:       primitive.NewObjectID(),
			Title:    "Test Title",
			Content:  "Test Content",
			Author:   "Test Author",
			AuthorID: "Test Author ID",
		})
		log.Println(blogPostID)
		require.NoError(t, err)
		require.NotNil(t, blogPostID)

		mt.ClearMockResponses()

	})
}

func TestFindBlogPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("FindBlogPost", func(mt *mtest.T) {
		blogPostCollection := mt.Coll
		blogPostDB := &BlogPostDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: blogPostCollection,
		}
		expectedBlogPost := models.BlogPost{
			ID:       primitive.NewObjectID(),
			Title:    "Test Title",
			Content:  "Test Content",
			Author:   "Test Author",
			AuthorID: "Test Author ID",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedBlogPost.ID},
			{Key: "title", Value: expectedBlogPost.Title},
			{Key: "content", Value: expectedBlogPost.Content},
			{Key: "author", Value: expectedBlogPost.Author},
			{Key: "authorID", Value: expectedBlogPost.AuthorID},
		}))

		blogPost, err := blogPostDB.FindBlogPost(bson.M{
			"ID": expectedBlogPost.ID,
		})

		log.Println(blogPost.ID)
		require.NoError(t, err)
		require.NotNil(t, blogPost)
		require.Equal(t, expectedBlogPost.ID, blogPost.ID)

		mt.ClearMockResponses()
	})
}

func TestFindAllBlogPosts(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("FindAllBlogPosts", func(mt *mtest.T) {
		blogPostCollection := mt.Coll
		blogPostDB := &BlogPostDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: blogPostCollection,
		}
		firstBlogPost := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "title", Value: "Test Title"},
			{Key: "content", Value: "Test Content"},
			{Key: "author", Value: "Test Author"},
			{Key: "authorID", Value: "Test Author ID"},
		})
		secondBlogPost := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "title", Value: "Test Title"},
			{Key: "content", Value: "Test Content"},
			{Key: "author", Value: "Test Author"},
			{Key: "authorID", Value: "Test Author ID"},
		})

		destroyCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(firstBlogPost, secondBlogPost, destroyCursors)

		blogPosts, err := blogPostDB.FindAllBlogPosts(bson.M{
			"authorID": "Test Author ID",
		})

		require.NoError(t, err)
		require.NotNil(t, blogPosts)
		require.Len(t, blogPosts, 2)
	})
}

func TestUpdateBlogPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("UpdateBlogPost", func(mt *mtest.T) {
		blogPostCollection := mt.Coll
		blogPostDB := &BlogPostDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: blogPostCollection,
		}
		expectedBlogPost := models.BlogPost{
			ID:       primitive.NewObjectID(),
			Title:    "Title",
			Content:  "Test Content",
			Author:   "Test Author",
			AuthorID: "Test Author ID",
		}
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: expectedBlogPost.ID},
				{Key: "title", Value: expectedBlogPost.Title},
				{Key: "content", Value: expectedBlogPost.Content},
				{Key: "author", Value: expectedBlogPost.Author},
				{Key: "authorID", Value: expectedBlogPost.AuthorID},
			}},
		})

		err := blogPostDB.UpdateBlogPost(bson.M{"ID": expectedBlogPost.ID}, bson.M{"$set": bson.M{"title": "New Title"}})
		require.NoError(t, err)
	})
}

func TestDeleteBlogPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("DeleteBlogPost", func(mt *mtest.T) {
		blogPostCollection := mt.Coll
		blogPostDB := &BlogPostDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: blogPostCollection,
		}
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})
		err := blogPostDB.DeleteBlogPost(bson.M{"ID": primitive.NewObjectID()})
		require.NoError(t, err)
	})
}

func TestDeleteAllBlogPosts(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("DeleteAllBlogPosts", func(mt *mtest.T) {
		blogPostCollection := mt.Coll
		blogPostDB := &BlogPostDB{
			MongoDB:    &MongoDB{client: mt.Client},
			Collection: blogPostCollection,
		}
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})
		err := blogPostDB.DeleteAllBlogPosts(bson.M{"authorID": "Test Author ID"})
		require.NoError(t, err)
	})
}
