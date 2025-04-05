package business

import (
	"context"
	"errors"

	"github.com/StephenJonesIT/CCMS/src/user-service/internal/models"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/repository"
	"github.com/google/uuid"
)

type FollowService interface {
	Follow(ctx context.Context, follow *models.Follows) error
	Unfollow(ctx context.Context, follow *models.Follows) error
	GetFollowers(ctx context.Context, followerID uuid.UUID) ([]models.Profile, error)
	GetFollowees(ctx context.Context, followerID uuid.UUID) ([]models.Profile, error)
}

type FollowServiceImpl struct {
	profileRepo repository.ProfileRepository
	followRepo repository.FollowRepository
}

func NewFollowService(profileRepo repository.ProfileRepository, followRepo repository.FollowRepository) *FollowServiceImpl{
	return &FollowServiceImpl{
		profileRepo: profileRepo,
		followRepo: followRepo,
	}
}

func (s *FollowServiceImpl)Follow(ctx context.Context, follow *models.Follows) error{
	if follow == nil {
		return errors.New("follow object cannot be nil")
	}

	if follow.FolloweeID == uuid.Nil || follow.FollowerID == uuid.Nil{
		return errors.New("FolloweeID or FollowerID is required")
	}

	if follow.FolloweeID == follow.FollowerID {
		return errors.New("unable to follow oneself")
	}

	if err := s.followRepo.Follow(ctx, follow); err != nil {
		return errors.New("Failed to create follow: "+ err.Error())
	}

	if err := s.profileRepo.IncrementFollowersCount(ctx,follow.FolloweeID); err != nil {
		s.followRepo.Unfollow(ctx, follow)
		return errors.New("Failed to increment count follower: "+err.Error())
	}

	if err := s.profileRepo.IncrementFolloweeCount(ctx,follow.FollowerID); err != nil {
		s.followRepo.Unfollow(ctx, follow)
		s.profileRepo.DecrementFollowersCount(ctx, follow.FolloweeID)
		return errors.New("Failed to increment count followee: "+err.Error())
	}
	return nil
}

func (s *FollowServiceImpl)Unfollow(ctx context.Context, follow *models.Follows) error {
	if follow == nil {
		return errors.New("follow object cannot be nil")
	}

	if follow.FolloweeID == uuid.Nil || follow.FollowerID == uuid.Nil{
		return errors.New("FolloweeID or FollowerID is required")
	}

	if follow.FolloweeID == follow.FollowerID {
		return errors.New("unable to unfollow oneself")
	}

	if err := s.followRepo.Unfollow(ctx, follow); err != nil {
		return errors.New("Failed to delete follow: "+ err.Error())
	}

	if err := s.profileRepo.DecrementFollowersCount(ctx,follow.FolloweeID); err != nil {
		s.followRepo.Follow(ctx, follow)
		return errors.New("Failed to decrement count follower: "+err.Error())
	}

	if err := s.profileRepo.DecrementFolloweeCount(ctx,follow.FollowerID); err != nil {
		s.followRepo.Follow(ctx, follow)
		s.profileRepo.IncrementFollowersCount(ctx, follow.FolloweeID)
		return errors.New("Failed to decrement count followee: "+err.Error())
	}
	return nil
}

func (s *FollowServiceImpl)GetFollowers(ctx context.Context, followID uuid.UUID) ([]models.Profile, error){
	if followID == uuid.Nil {
		return nil, errors.New("FollowID is required")
	}

	result, err := s.profileRepo.GetFollowers(ctx, followID)
	if  err != nil {
		return nil, errors.New("Failed to retieve followers: "+err.Error())
	}

	return result, nil
}

func (s *FollowServiceImpl)GetFollowees(ctx context.Context, followID uuid.UUID) ([]models.Profile, error){
	if followID == uuid.Nil {
		return nil, errors.New("FollowID is required")
	}

	result, err := s.profileRepo.GetFollowees(ctx, followID)
	if  err != nil {
		return nil, errors.New("Failed to retieve followees: "+err.Error())
	}

	return result, nil
}