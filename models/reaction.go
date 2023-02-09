package models

type Reaction struct {
	Id       int
	PostId   int
	Reaction string
	Username string
	Object   string
}
