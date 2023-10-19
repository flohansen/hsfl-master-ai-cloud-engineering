package service

import (
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/models"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/repository"
)

// PostServiceImpl handles business logic for Post
type PostServiceImpl struct {
	PostRepository repository.PostRepository
}

// NewPostService creates a new PostServiceImpl
func NewPostService(repo repository.PostRepository) *PostServiceImpl {
	return &PostServiceImpl{PostRepository: repo}
}

// Create a new post
func (s *PostServiceImpl) Create(post *models.Post) {
	s.PostRepository.Create(post)
}

// GetAll Get all posts
func (s *PostServiceImpl) GetAll() []models.Post {
	return s.PostRepository.FindAll()
}

// GetByID Get a post by ID
func (s *PostServiceImpl) GetByID(id uint) models.Post {
	return s.PostRepository.FindByID(id)
}

// Update a post
func (s *PostServiceImpl) Update(post *models.Post) {
	s.PostRepository.Update(post)
}

// Delete a post
func (s *PostServiceImpl) Delete(post *models.Post) {
	s.PostRepository.Delete(post)
}
