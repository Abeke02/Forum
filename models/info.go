package models

type Info struct {
	User       User
	ThatUser   User
	Posts      []Post
	PostsLike  []Post
	Post       Post
	Comments   []Comment
	Categories []string
	Category   string
}
