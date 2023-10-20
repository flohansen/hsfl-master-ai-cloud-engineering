package service

import (
	mocks "github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/_mocks"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/models"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	postService := NewPostService(mockRepo)
	post := models.Post{Title: "Test Post", Content: "Test Content"}

	mockRepo.EXPECT().Create(&post)

	postService.Create(&post)
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	postService := NewPostService(mockRepo)

	mockRepo.EXPECT().FindAll().Return([]models.Post{})

	postService.GetAll()
}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	postService := NewPostService(mockRepo)

	mockRepo.EXPECT().FindByID(uint(1)).Return(models.Post{})

	postService.GetByID(uint(1))
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	postService := NewPostService(mockRepo)
	post := models.Post{Title: "Test Post", Content: "Test Content"}

	mockRepo.EXPECT().Update(&post)

	postService.Update(&post)
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	postService := NewPostService(mockRepo)
	post := models.Post{Title: "Test Post", Content: "Test Content"}

	mockRepo.EXPECT().Delete(&post)

	postService.Delete(&post)
}
