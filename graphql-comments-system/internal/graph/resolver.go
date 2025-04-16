package graph

import (
	"sync"

	"github.com/yyyppd/graphql-comments-system/internal/comment"
	"github.com/yyyppd/graphql-comments-system/internal/graph/model"
	"github.com/yyyppd/graphql-comments-system/internal/post"
)

type Resolver struct {
	PostRepo           post.PostRepository
	CommentRepo        comment.CommentRepository
	mu                 sync.Mutex
	commentSubscribers map[string][]chan *model.Comment
}

func NewResolver(postRepo post.PostRepository, commentRepo comment.CommentRepository) *Resolver {
	return &Resolver{
		PostRepo:           postRepo,
		CommentRepo:        commentRepo,
		commentSubscribers: make(map[string][]chan *model.Comment),
	}
}
