package service

import (
	"main.go/internal/storage"
	"main.go/models"
)

type Comment interface {
	CreateComment(id int, user models.User, comment string) error
	GetCommentsByIdPost(id int) ([]models.Comment, error)
}

type CommentService struct {
	storage storage.Comment
}

func newCommentService(storage storage.Comment) *CommentService {
	return &CommentService{
		storage: storage,
	}
}

func (c *CommentService) CreateComment(id int, user models.User, comment string) error {
	if err := c.storage.CreateComment(id, user, comment); err != nil {
		return err
	}
	return nil
}

func (c *CommentService) GetCommentsByIdPost(id int) ([]models.Comment, error) {
	comments, err := c.storage.GetCommentsByIdPost(id)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
