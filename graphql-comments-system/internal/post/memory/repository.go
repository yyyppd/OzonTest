package memory

import (
	"context"
	"sync"

	"github.com/yyyppd/graphql-comments-system/internal/post"
)

type InMemoryPostRepo struct {
	mu     sync.RWMutex
	posts  map[int]*post.Post
	nextID int
}

func New() *InMemoryPostRepo {
	return &InMemoryPostRepo{
		posts:  make(map[int]*post.Post),
		nextID: 1,
	}
}

func (r *InMemoryPostRepo) Create(ctx context.Context, title, content string, allowComments bool) (*post.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	p := &post.Post{
		ID:            id,
		Title:         title,
		Content:       content,
		AllowComments: allowComments,
	}
	r.posts[id] = p
	return p, nil
}

func (r *InMemoryPostRepo) GetAll(ctx context.Context) ([]*post.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*post.Post
	for _, p := range r.posts {
		result = append(result, p)
	}
	return result, nil
}

func (r *InMemoryPostRepo) GetByID(ctx context.Context, id int) (*post.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.posts[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}
