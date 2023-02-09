package service

import "main.go/internal/storage"

type Service struct {
	Auth
	Post
	User
	Comment
	Reaction
}

func NewService(storages *storage.Storage) *Service {
	return &Service{
		Auth:     newAuthService(storages.Auth),
		Post:     newPostService(storages.Post),
		User:     newUserService(storages.User),
		Comment:  newCommentService(storages.Comment),
		Reaction: newReactionService(storages.Reaction),
	}
}
