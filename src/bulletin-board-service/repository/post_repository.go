package repository

import "github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/models"

type PostRepository interface {
	Create(post *models.Post)
	FindAll() []models.Post
	FindByID(id uint) models.Post
	Update(post *models.Post)
	Delete(post *models.Post)
}
