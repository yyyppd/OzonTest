// internal/comment/memory/repository_test.go
package memory

import (
	"context"
	"testing"
)

func TestCreateAndGetByPostID(t *testing.T) {
	repo := New()
	ctx := context.Background()

	c, err := repo.Create(ctx, 1, nil, "Test comment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	comments, err := repo.GetByPostID(ctx, 1, 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}

	if comments[0].ID != c.ID {
		t.Errorf("expected comment ID %d, got %d", c.ID, comments[0].ID)
	}
}
