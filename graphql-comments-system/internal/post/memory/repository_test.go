package memory

import (
	"context"
	"testing"
)

func TestCreateAndGetByID(t *testing.T) {
	repo := New()
	ctx := context.Background()

	created, err := repo.Create(ctx, "Memory Title", "Memory Content", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == nil || got.Title != "Memory Title" {
		t.Errorf("expected title %q, got %+v", "Memory Title", got)
	}
}

func TestGetAll(t *testing.T) {
	repo := New()
	ctx := context.Background()

	_, _ = repo.Create(ctx, "Post 1", "Content 1", true)
	_, _ = repo.Create(ctx, "Post 2", "Content 2", false)

	all, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(all) != 2 {
		t.Errorf("expected 2 posts, got %d", len(all))
	}
}
