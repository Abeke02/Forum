package models

import "time"

type Post struct {
	Id          int
	Creator     string
	Title       string
	Description string
	Category    []string
	CreatedAt   time.Time
	Likes       int
	Dislikes    int
}
