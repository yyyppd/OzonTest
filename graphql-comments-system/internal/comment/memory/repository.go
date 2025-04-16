package memory

import (
	"context"
	"sync"
	"time"

	"github.com/yyyppd/graphql-comments-system/internal/comment"
)

type InMemoryCommentRepo struct {
	mu       sync.RWMutex
	comments map[int]*comment.Comment
	nextID   int
}

func New() *InMemoryCommentRepo {
	return &InMemoryCommentRepo{
		comments: make(map[int]*comment.Comment),
		nextID:   1,
	}
}

func (r *InMemoryCommentRepo) Create(ctx context.Context, postID int, parentID *int, content string) (*comment.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	newComment := &comment.Comment{
		ID:        id,
		PostID:    postID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if parentID != nil {
		newComment.ParentID.Int64 = int64(*parentID)
		newComment.ParentID.Valid = true
	}

	r.comments[id] = newComment

	return newComment, nil
}

func (r *InMemoryCommentRepo) GetByPostID(ctx context.Context, postID, limit, offset int) ([]*comment.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*comment.Comment
	count := 0

	for _, c := range r.comments {
		if c.PostID == postID && !c.ParentID.Valid {
			if count >= offset && len(result) < limit {
				result = append(result, c)
			}
			count++
		}
	}

	return result, nil
}

func (r *InMemoryCommentRepo) GetChildren(ctx context.Context, parentID int) ([]*comment.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*comment.Comment
	for _, c := range r.comments {
		if c.ParentID.Valid && int(c.ParentID.Int64) == parentID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryCommentRepo) GetChildrenBatch(ctx context.Context, parentIDs []int) ([]*comment.Comment, error) {
	var result []*comment.Comment
	parentSet := make(map[int]struct{}, len(parentIDs))
	for _, id := range parentIDs {
		parentSet[id] = struct{}{}
	}

	for _, c := range r.comments {
		if c.ParentID.Valid {
			if _, ok := parentSet[int(c.ParentID.Int64)]; ok {
				result = append(result, c)
			}
		}
	}

	return result, nil
}
