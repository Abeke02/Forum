package storage

import (
	"database/sql"
	"fmt"

	"main.go/models"
)

type Comment interface {
	CreateComment(id int, user models.User, comment string) error
	GetCommentsByIdPost(id int) ([]models.Comment, error)
}

type CommentStorage struct {
	db *sql.DB
}

func newCommentStorage(db *sql.DB) *CommentStorage {
	return &CommentStorage{
		db: db,
	}
}

func (c *CommentStorage) CreateComment(id int, user models.User, comment string) error {
	query := `INSERT INTO comment(id_post, creator, comment) VALUES($1, $2, $3);`
	if _, err := c.db.Exec(query, id, user.Username, comment); err != nil {
		return fmt.Errorf("storage: add comment: %w", err)
	}
	return nil
}

func (c *CommentStorage) GetCommentsByIdPost(id int) ([]models.Comment, error) {
	comments := []models.Comment{}
	query := `SELECT id, creator, comment,created_at FROM comment WHERE id_post=$1;`
	rows, err := c.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("storage: comment by id post: %w", err)
	}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.Id, &comment.Creator, &comment.Text, &comment.Created_at); err != nil {
			return nil, fmt.Errorf("storage: comment by id post: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, err
}
