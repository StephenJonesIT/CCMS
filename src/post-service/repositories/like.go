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

type LikeRepository interface {
	ToggleLike(ctx context.Context, userID uuid.UUID, targetID primitive.ObjectID, targetType string) (bool, error)
	GetPostAuthorID(ctx context.Context, postID primitive.ObjectID) (uuid.UUID, error)
	GetCommentAuthorID(ctx context.Context, commentID primitive.ObjectID) (uuid.UUID, error)
	DeleteLikesByTarget(ctx context.Context, targetID primitive.ObjectID, targetType string) error
}

type LikeRepositoryImpl struct {
	db *mongo.Database
}

func NewLikeRepositoryImpl(db *mongo.Client) *LikeRepositoryImpl {
	return &LikeRepositoryImpl{db: config.GetDatabase(db)}
}

func (r *LikeRepositoryImpl) ToggleLike(ctx context.Context, userID uuid.UUID, targetID primitive.ObjectID, targetType string) (bool, error) {
	// Kiểm tra nếu target là comment và status không phải "approved"
	if targetType == "comment" {
		var comment struct {
			Status string `bson:"status"`
		}
		err := r.db.Collection("comments").FindOne(ctx, bson.M{"_id": targetID}).Decode(&comment)
		if err != nil {
			return false, err
		}
		if comment.Status != "approved" {
			return false, errors.New("cannot like unapproved comment")
		}
	}

	filter := bson.M{
		"user_id":     userID,
		"target_id":   targetID,
		"target_type": targetType,
	}

	var existingLike models.Like
	err := r.db.Collection("likes").FindOne(ctx, filter).Decode(&existingLike)

	if err == nil {
		// Like exists - unlike
		_, err := r.db.Collection("likes").DeleteOne(ctx, filter)
		if err != nil {
			return false, err
		}
		return false, r.decrementLikeCount(ctx, targetID, targetType)
	} else if err != mongo.ErrNoDocuments {
		// Real error occurred
		return false, err
	}

	// Like doesn't exist - create new like
	newLike := models.Like{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
		CreatedAt:  time.Now(),
	}

	_, err = r.db.Collection("likes").InsertOne(ctx, newLike)
	if err != nil {
		return false, err
	}

	// Increment like count
	if err := r.incrementLikeCount(ctx, targetID, targetType); err != nil {
		return false, err
	}

	return true, nil
}

func (r *LikeRepositoryImpl) incrementLikeCount(ctx context.Context, targetID primitive.ObjectID, targetType string) error {
	var collection string
	if targetType == "post" {
		collection = "posts"
	} else if targetType == "comment" {
		collection = "comments"
	} else {
		return errors.New("invalid target type")
	}

	filter := bson.M{"_id": targetID}
	update := bson.M{"$inc": bson.M{"like_count": 1}}

	_, err := r.db.Collection(collection).UpdateOne(ctx, filter, update)
	return err
}

func (r *LikeRepositoryImpl) decrementLikeCount(ctx context.Context, targetID primitive.ObjectID, targetType string) error {
	var collection string
	if targetType == "post" {
		collection = "posts"
	} else if targetType == "comment" {
		collection = "comments"
	} else {
		return errors.New("invalid target type")
	}

	filter := bson.M{"_id": targetID}
	update := bson.M{"$inc": bson.M{"like_count": -1}}

	_, err := r.db.Collection(collection).UpdateOne(ctx, filter, update)
	return err
}

func (r *LikeRepositoryImpl) GetPostAuthorID(ctx context.Context, postID primitive.ObjectID) (uuid.UUID, error) {
	var post struct {
		AuthorID uuid.UUID `bson:"author_id"`
	}
	err := r.db.Collection("posts").FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		return uuid.Nil, err
	}
	return post.AuthorID, nil
}

func (r *LikeRepositoryImpl) GetCommentAuthorID(ctx context.Context, commentID primitive.ObjectID) (uuid.UUID, error) {
	var comment struct {
		UserID uuid.UUID `bson:"user_id"`
	}
	err := r.db.Collection("comments").FindOne(ctx, bson.M{"_id": commentID}).Decode(&comment)
	if err != nil {
		return uuid.Nil, err
	}
	return comment.UserID, nil
}

// like_repository.go
func (r *LikeRepositoryImpl) DeleteLikesByTarget(ctx context.Context, targetID primitive.ObjectID, targetType string) error {
    filter := bson.M{
        "target_id":   targetID,
        "target_type": targetType,
    }
    _, err := r.db.Collection("likes").DeleteMany(ctx, filter)
    return err
}