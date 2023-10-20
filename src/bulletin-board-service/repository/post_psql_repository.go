package repository

import (
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/models"
	"gorm.io/gorm"
)

// PostPsqlRepository handles database operations for Post
type PostPsqlRepository struct {
	DB *gorm.DB
}

// NewPostPsqlRepository creates a new PostPsqlRepository
func NewPostPsqlRepository(db *gorm.DB) *PostPsqlRepository {
	return &PostPsqlRepository{DB: db}
}

// Create a new post
func (r *PostPsqlRepository) Create(post *models.Post) {
	r.DB.Create(post)
}

// FindAll posts
func (r *PostPsqlRepository) FindAll() []models.Post {
	var posts []models.Post
	r.DB.Find(&posts)
	return posts
}

// FindByID finds a post by ID
func (r *PostPsqlRepository) FindByID(id uint) models.Post {
	var post models.Post
	r.DB.First(&post, id)
	return post
}

// Update a post
func (r *PostPsqlRepository) Update(post *models.Post) {
	r.DB.Save(post)
}

// Delete a post
func (r *PostPsqlRepository) Delete(post *models.Post) {
	r.DB.Delete(post)
}
