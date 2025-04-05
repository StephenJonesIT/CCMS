package business

import (
	"context"
	"github.com/StephenJonesIT/CCMS/src/post-service/repositories"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeService interface {
	ToggleLike(ctx context.Context, userID uuid.UUID, targetID primitive.ObjectID, targetType string) (bool, error)
}

type LikeServiceImpl struct {
	likeRepo repositories.LikeRepository
}

func NewLikeService(likeRepo repositories.LikeRepository) *LikeServiceImpl {
	return &LikeServiceImpl{likeRepo: likeRepo}
}

func (s *LikeServiceImpl) ToggleLike(ctx context.Context, userID uuid.UUID, targetID primitive.ObjectID, targetType string) (bool, error) {
	isLiked, err := s.likeRepo.ToggleLike(ctx, userID, targetID, targetType)
	if err != nil {
		return false, err
	}

	// Gửi thông báo nếu là like mới
	// if isLiked {
	// 	var authorID uuid.UUID
	// 	var message string

	// 	if targetType == "post" {
	// 		authorID, err = s.likeRepo.GetPostAuthorID(ctx, targetID)
	// 		message = "liked your post"
	// 	} else {
	// 		authorID, err = s.likeRepo.GetCommentAuthorID(ctx, targetID)
	// 		message = "liked your comment"
	// 	}

	// 	if err != nil {
	// 		log.Printf("failed to get author ID: %v", err)
	// 	} else if authorID != userID {
	// 		notification := &Notification{
	// 			UserID:     authorID,
	// 			ActionUser: userID,
	// 			TargetID:   targetID,
	// 			TargetType: targetType,
	// 			Message:    message,
	// 			IsRead:     false,
	// 			CreatedAt:  time.Now(),
	// 		}

	// 		if err := s.notifRepo.CreateNotification(ctx, notification); err != nil {
	// 			log.Printf("failed to create notification: %v", err)
	// 		}
	// 	}
	// }

	return isLiked, nil
}