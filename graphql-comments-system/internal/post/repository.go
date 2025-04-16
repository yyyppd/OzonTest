package post

import "context"

type Post struct {
	ID            int
	Title         string
	Content       string
	AllowComments bool
}

type PostRepository interface {
	GetAll(ctx context.Context) ([]*Post, error)
	GetByID(ctx context.Context, id int) (*Post, error)
	Create(ctx context.Context, title, content string, allowComments bool) (*Post, error)
}
