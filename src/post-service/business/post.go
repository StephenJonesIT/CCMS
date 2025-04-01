package business

import (
	"errors"
	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/post-service/models"
	"github.com/StephenJonesIT/CCMS/src/post-service/repositories"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
	CreatePost(post *models.Post) (*models.Post,error)
	GetPostByID(id string) (*models.Post, error)
	UpdatePost(id string, post *models.Post) error
	DeletePost(id string) error
	GetPostsByUser(userID string, page int, limit int) ([]models.Post, error)
	GetPostsByStatus(page common.Paging, status string) ([]models.Post, error)
}

type PostServiceImpl struct {
	postRepo repositories.PostRepository
}

func NewPostService(postRepo repositories.PostRepository) *PostServiceImpl {
	return &PostServiceImpl{
		postRepo: postRepo,
	}
}

func(p *PostServiceImpl) CreatePost(post *models.Post)(*models.Post, error){
	if post.AuthorID == uuid.Nil {
		return  nil, errors.New("AuthorID is required")
	}

	if post.Status == "" {
		return  nil ,errors.New("status is required")
	}
    
	return p.postRepo.CreatePost(post)
}

func(p *PostServiceImpl) GetPostByID(id string)(*models.Post, error){
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	post, err := p.postRepo.GetPostByID(objectID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func(p *PostServiceImpl) UpdatePost(id string, post *models.Post) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = p.postRepo.UpdatePost(objectID, post)
	if err != nil {
		return err
	}

	return nil
}

func(p *PostServiceImpl) DeletePost(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = p.postRepo.DeletePost(objectId)
	if err != nil {
		return err
	}

	return nil
}

func(p *PostServiceImpl) GetPostsByUser(userID string, page int, limit int) ([]models.Post, error){
	id , err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Validate page and limit
	if page <= 0 {
		page = 1
	}
	
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// Fetch posts from the repository
	posts, err := p.postRepo.GetPostsByUser(id, page, limit)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, errors.New("no posts found")
	}

	return posts, nil
}

func(p *PostServiceImpl) GetPostsByStatus(page common.Paging, status string) ([]models.Post, error){
	return nil, nil
}