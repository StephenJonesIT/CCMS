package repositories

import (
	"context"
	"errors"
	"time"
	"github.com/StephenJonesIT/CCMS/src/post-service/config"
	"github.com/StephenJonesIT/CCMS/src/post-service/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentByID(ctx context.Context, id primitive.ObjectID) (*models.Comment, error)
	UpdateComment(ctx context.Context, id primitive.ObjectID, userID uuid.UUID, content string) error
	DeleteComment(ctx context.Context, id primitive.ObjectID, userID uuid.UUID) error
	AddReply(ctx context.Context, commentID primitive.ObjectID, reply *models.Reply) error
	UpdateCommentStatus(ctx context.Context, id primitive.ObjectID, status string) error
	GetCommentsByPostID(ctx context.Context, postID primitive.ObjectID, status string) ([]models.Comment, error)
}

type CommentRepositoryImpl struct {
	collection *mongo.Collection
}

func NewCommentRepositoryImpl(db *mongo.Client) *CommentRepositoryImpl {
	return &CommentRepositoryImpl{
		collection: config.GetCollection(db, "comments"),
	}
}

func (r *CommentRepositoryImpl) CreateComment(ctx context.Context, comment *models.Comment) error {
	comment.CreatedAt = time.Now()
	comment.Status = "pending" // Mặc định là chờ kiểm duyệt
	_, err := r.collection.InsertOne(ctx, comment)
	return err
}

func (r *CommentRepositoryImpl) GetCommentByID(ctx context.Context, id primitive.ObjectID) (*models.Comment, error) {
	var comment models.Comment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepositoryImpl) UpdateComment(ctx context.Context, id primitive.ObjectID, userID uuid.UUID, content string) error {
	filter := bson.M{"_id": id, "user_id": userID}
	update := bson.M{"$set": bson.M{
		"content":    content,
		"updated_at": time.Now(),
	}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("comment not found or not owned by user")
	}
	return nil
}

func (r *CommentRepositoryImpl) DeleteComment(ctx context.Context, id primitive.ObjectID, userID uuid.UUID) error {
	filter := bson.M{"_id": id, "user_id": userID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("comment not found or not owned by user")
	}
	return nil
}

func (r *CommentRepositoryImpl) AddReply(ctx context.Context, commentID primitive.ObjectID, reply *models.Reply) error {
	reply.CreatedAt = time.Now()
	filter := bson.M{"_id": commentID}
	update := bson.M{"$push": bson.M{"replies": reply}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *CommentRepositoryImpl) UpdateCommentStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *CommentRepositoryImpl) GetCommentsByPostID(ctx context.Context, postID primitive.ObjectID, status string) ([]models.Comment, error) {
	filter := bson.M{"post_id": postID}
	if status != "" {
		filter["status"] = status
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}