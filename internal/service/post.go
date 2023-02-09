package service

import (
	"strings"

	"main.go/internal/storage"
	"main.go/models"
)

type Post interface {
	//	CRUD
	CreatePost(post models.Post, user models.User) error
	// UpdatePost(postId int, post models.Post, user models.User) error
	// DeletePost(postId int, user models.User) error
	//	GET
	GetAllPosts(filter string) ([]models.Post, error)
	GetPostById(id int) (models.Post, error)
	GetPostsByUsername(username string) ([]models.Post, error)
	GetAllCategories() ([]string, error)
	UpdateCountsReactions(object string, likes int, dislikes int, postId int) error
	GetIDPostsByUsername(username string) ([]models.Reaction, error)
}

type PostService struct {
	storage storage.Post
}

func newPostService(storage storage.Post) *PostService {
	return &PostService{
		storage: storage,
	}
}

func (p *PostService) CreatePost(post models.Post, user models.User) error {
	post.Category = strings.Fields(post.Category[0])

	if err := p.storage.CreatePost(user.Username, post); err != nil {
		return err
	}
	return nil
}

func (p *PostService) GetPostById(id int) (models.Post, error) {
	post, err := p.storage.GetPostById(id)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (p *PostService) GetAllPosts(filter string) ([]models.Post, error) {
	return p.storage.GetAllPosts(filter)
}

func (p *PostService) GetPostsByUsername(username string) ([]models.Post, error) {
	posts, err := p.storage.GetPostsByUsername(username)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostService) GetIDPostsByUsername(username string) ([]models.Reaction, error) {
	postsId, err := p.storage.GetIDPostsByUsername(username)
	if err != nil {
		return nil, err
	}

	return postsId, nil
}

func (p *PostService) UpdateCountsReactions(object string, likes int, dislikes int, postId int) error {
	if err := p.storage.UpdateCountsReactions(object, likes, dislikes, postId); err != nil {
		return err
	}
	return nil
}

func (p *PostService) GetAllCategories() ([]string, error) {
	var categories []string
	categories, err := p.storage.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}
