package service

import "github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/models"

type PostService interface {
	Create(post *models.Post)
	GetAll() []models.Post
	GetByID(id uint) models.Post
	Update(post *models.Post)
	Delete(post *models.Post)
}
