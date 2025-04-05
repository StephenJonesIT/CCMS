package repositories

import (
	"context"
	"time"

	"github.com/StephenJonesIT/CCMS/src/post-service/config"
	"github.com/StephenJonesIT/CCMS/src/post-service/models"
	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type PostRepository interface {
	CreatePost(post *models.Post) (*models.Post, error)
	GetPostByID(id primitive.ObjectID) (*models.Post, error)
	UpdatePost(id primitive.ObjectID, post *models.Post) error
	DeletePost(id primitive.ObjectID,) error
	GetPostsByUser(userID uuid.UUID, page int, limit int) ([]models.Post, error)
	GetPostsByStatus(page common.Paging, status string) ([]models.Post,error)
	DecrementCommentCount(ctx context.Context, postID primitive.ObjectID) error
	IncrementCommentCount(ctx context.Context ,postID primitive.ObjectID) error
}

type PostRepositoryImpl struct {
	collection *mongo.Collection
}

func NewPostRepositoryImpl(db *mongo.Client) *PostRepositoryImpl {
	return &PostRepositoryImpl{
		collection: config.GetCollection(db, "posts"),
	}
}

func (r *PostRepositoryImpl) CreatePost(post *models.Post) (*models.Post, error) {

	result, err := r.collection.InsertOne(context.Background(), post)
	if err != nil {
		return nil, err
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	return post, nil
}

func (r *PostRepositoryImpl) GetPostByID(id primitive.ObjectID) (*models.Post, error) {

	var post models.Post
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepositoryImpl) UpdatePost(id primitive.ObjectID, post *models.Post) error {
	post.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": post},
	)

	return err
}

func (r *PostRepositoryImpl) DeletePost(id primitive.ObjectID,) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *PostRepositoryImpl) GetPostsByUser(userID uuid.UUID, page int, limit int) ([]models.Post, error) {
	skip := (page - 1) * limit

	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"author_id": userID},
		options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []models.Post
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepositoryImpl) GetPostsByStatus(page common.Paging, status string) ([]models.Post,error) {
	skip := (page.Page-1) *page.Limit
	limit := page.Limit
	filter := bson.M{"status": status}

	cussor, err := r.collection.Find(
		context.Background(), 
		filter, 
		options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1}),
	)

	if err != nil {
		return nil, err
	}
	defer cussor.Close(context.Background())
	var posts []models.Post
	if err = cussor.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// In PostRepositoryImpl add:
func (r *PostRepositoryImpl) IncrementCommentCount(ctx context.Context, postID primitive.ObjectID) error {
    filter := bson.M{"_id": postID}
    update := bson.M{
        "$inc": bson.M{"comment_count": 1},
        "$set": bson.M{"updated_at": time.Now()},
    }

    _, err := r.collection.UpdateOne(
        ctx,
        filter,
        update,
    )
    
    return err
}

func (r *PostRepositoryImpl) DecrementCommentCount(ctx context.Context, postID primitive.ObjectID) error {
    filter := bson.M{
        "_id": postID,
        "comment_count": bson.M{"$gt": 0}, // Only decrement if count > 0
    }
    
    update := bson.M{
        "$inc": bson.M{"comment_count": -1},
        "$set": bson.M{"updated_at": time.Now()},
    }

    _, err := r.collection.UpdateOne(
        ctx,
        filter,
        update,
    )
    
    return err
}