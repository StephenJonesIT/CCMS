package business

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/StephenJonesIT/CCMS/src/post-service/models"
	"github.com/StephenJonesIT/CCMS/src/post-service/repositories"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService interface {
	CreateComment(ctx context.Context, postID primitive.ObjectID, userID uuid.UUID, content string) (*models.Comment, error)
	UpdateComment(ctx context.Context, commentID primitive.ObjectID, userID uuid.UUID, content string) error
	DeleteComment(ctx context.Context, commentID primitive.ObjectID, userID uuid.UUID) error
	AddReply(ctx context.Context, commentID primitive.ObjectID, userID uuid.UUID, content string) error
	GetCommentsByPostID(ctx context.Context, postID primitive.ObjectID, showPending bool) ([]models.Comment, error)
	ModerateComment(ctx context.Context, commentID primitive.ObjectID, status string) error
}

type CommentServiceImpl struct {
    commentRepo repositories.CommentRepository
    postRepo    repositories.PostRepository
    likeRepo    repositories.LikeRepository
}

func NewCommentService(
	commentRepo repositories.CommentRepository, 
	postRepo repositories.PostRepository, 
	likeRepo repositories.LikeRepository,
	) *CommentServiceImpl {
	return &CommentServiceImpl{
		commentRepo: commentRepo,
		postRepo: postRepo,
	}
}

func (s *CommentServiceImpl) CreateComment(ctx context.Context, postID primitive.ObjectID, userID uuid.UUID, content string) (*models.Comment, error) {
	comment := &models.Comment{
		ID: primitive.NewObjectID(),		
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		Replies:   []models.Reply{},
		LikeCount: 0,
	}

	if err := s.commentRepo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentServiceImpl) UpdateComment(ctx context.Context, commentID primitive.ObjectID, userID uuid.UUID, content string) error {
	return s.commentRepo.UpdateComment(ctx, commentID, userID, content)
}

func (s *CommentServiceImpl) DeleteComment(ctx context.Context, commentID primitive.ObjectID, userID uuid.UUID) error {
    // Lấy thông tin comment trước khi xóa
    comment, err := s.commentRepo.GetCommentByID(ctx, commentID)
    if err != nil {
        return fmt.Errorf("failed to get comment: %w", err)
    }

    // Kiểm tra quyền sở hữu (trừ khi là admin/moderator)
    isAdmin := false // Giả sử lấy từ context hoặc middleware
    if comment.UserID != userID && !isAdmin {
        return errors.New("unauthorized: you can only delete your own comments")
    }

    // Xóa comment
    if err := s.commentRepo.DeleteComment(ctx, commentID, userID); err != nil {
        return fmt.Errorf("failed to delete comment: %w", err)
    }

    // Nếu comment đã được approved, giảm comment count trong post
    if comment.Status == "approved" {
        if err := s.postRepo.DecrementCommentCount(ctx, comment.PostID); err != nil {
            log.Printf("warning: failed to decrement comment count for post %s: %v", comment.PostID.Hex(), err)
            // Không return error vì comment đã xóa thành công
        }
    }

    // Xóa tất cả like liên quan đến comment này
    if err := s.likeRepo.DeleteLikesByTarget(ctx, commentID, "comment"); err != nil {
        log.Printf("warning: failed to delete likes for comment %s: %v", commentID.Hex(), err)
    }

    return nil
}

func (s *CommentServiceImpl) AddReply(ctx context.Context, commentID primitive.ObjectID, userID uuid.UUID, content string) error {
	reply := &models.Reply{
		UserID:    userID.String(),
		Content:   content,
	}

	return s.commentRepo.AddReply(ctx, commentID, reply)
}

func (s *CommentServiceImpl) GetCommentsByPostID(ctx context.Context, postID primitive.ObjectID, showPending bool) ([]models.Comment, error) {
	status := "approved"
	if showPending {
		status = ""
	}
	return s.commentRepo.GetCommentsByPostID(ctx, postID, status)
}

func (s *CommentServiceImpl) ModerateComment(ctx context.Context, commentID primitive.ObjectID, status string) error {
    if status != "approved" && status != "rejected" {
        return errors.New("invalid status")
    }

    // Lấy thông tin comment hiện tại để kiểm tra trạng thái cũ
    currentComment, err := s.commentRepo.GetCommentByID(ctx, commentID)
    if err != nil {
        return fmt.Errorf("failed to get comment: %w", err)
    }

    // Cập nhật trạng thái mới
    if err := s.commentRepo.UpdateCommentStatus(ctx, commentID, status); err != nil {
        return err
    }

    // Nếu chuyển từ pending sang approved, tăng comment count
    if currentComment.Status == "pending" && status == "approved" {
        if err := s.postRepo.IncrementCommentCount(ctx, currentComment.PostID); err != nil {
            log.Printf("failed to increment comment count for post %s: %v", currentComment.PostID.Hex(), err)
            // Không return error ở đây vì moderate đã thành công
        }
    } else if currentComment.Status == "approved" && status == "rejected" {
        if err := s.postRepo.DecrementCommentCount(ctx, currentComment.PostID); err != nil {
            log.Printf("failed to decrement comment count for post %s: %v", currentComment.PostID.Hex(), err)
        }
    }

    return nil
}