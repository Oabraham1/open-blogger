package db

import (
	"github.com/Oabraham1/open-blogger/server/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogPostDB struct {
	MongoDB 	*MongoDB
	Collection 	*mongo.Collection
}

func NewBlogPostDB() *BlogPostDB {
	mongoDB := ConnectToMongo()
	defer mongoDB.Disconnect()
	return &BlogPostDB{
		MongoDB: mongoDB,
	}
}

func (blogPostDB *BlogPostDB) CreateNewBlogPost(post models.BlogPost) (primitive.ObjectID, error) {
	blogPost, err := blogPostDB.MongoDB.InsertSingleDocument(blogPostDB.Collection, post)
	if err != nil {
		return primitive.NewObjectID(), err
	}
	var blogPostModel models.BlogPost
	blogPostModel.ID = blogPost.InsertedID.(primitive.ObjectID)
	return blogPostModel.ID, nil
}

func (blogPostDB *BlogPostDB) FindBlogPost(filter bson.M) (*models.BlogPost, error) {
	blogPost := blogPostDB.MongoDB.FindSingleDocument(blogPostDB.Collection, filter)
	if blogPost.Err() != nil {
		return nil, blogPost.Err()
	}

	var blogPostModel models.BlogPost
	err := blogPost.Decode(&blogPostModel)
	if err != nil {
		return nil, err
	}
	return &blogPostModel, nil
}

func (blogPostDB *BlogPostDB) FindAllBlogPosts(filter bson.M) ([]*models.BlogPost, error) {
	curr, err := blogPostDB.MongoDB.FindMultipleDocuments(blogPostDB.Collection, filter)
	if err != nil {
		return nil, err
	}

	defer curr.Close(blogPostDB.MongoDB.ctx)
	
	blogPosts := make([]*models.BlogPost, 0)

	for curr.Next(blogPostDB.MongoDB.ctx) {
		var blogPost models.BlogPost
		err := curr.Decode(&blogPost)
		if err != nil {
			return nil, err
		}
		blogPosts = append(blogPosts, &blogPost)
	}

	if err := curr.Err(); err != nil {
		return nil, err
	}

	return blogPosts, nil
}

func (blogPostDB *BlogPostDB) UpdateBlogPost(filter bson.M, update bson.M) error {
	_, err := blogPostDB.MongoDB.UpdateSingleDocument(blogPostDB.Collection, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (blogPostDB *BlogPostDB) DeleteBlogPost(filter bson.M) error {
	_, err := blogPostDB.MongoDB.DeleteSingleDocument(blogPostDB.Collection, filter)
	if err != nil {
		return err
	}
	return nil
}

func (blogPostDB *BlogPostDB) DeleteAllBlogPosts(filter bson.M) error {
	_, err := blogPostDB.MongoDB.DeleteMultipleDocuments(blogPostDB.Collection, filter)
	if err != nil {
		return err
	}
	return nil
}