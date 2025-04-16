package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := New(db)
	ctx := context.Background()

	postID := 1
	parentID := 2
	content := "Test comment"

	mock.ExpectQuery(`INSERT INTO comments`).
		WithArgs(postID, sql.NullInt64{Int64: int64(parentID), Valid: true}, content).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	c, err := repo.Create(ctx, postID, &parentID, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.ID != 123 {
		t.Errorf("expected ID 123, got %d", c.ID)
	}
	if c.Content != content {
		t.Errorf("expected content %q, got %q", content, c.Content)
	}
}

func TestGetChildren(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := New(db)
	ctx := context.Background()

	parentID := 42

	mock.ExpectQuery(`SELECT id, post_id, parent_id, content, created_at FROM comments WHERE parent_id = \$1`).
		WithArgs(parentID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "parent_id", "content", "created_at"}).
			AddRow(1, 1, parentID, "child 1", time.Now()).
			AddRow(2, 1, parentID, "child 2", time.Now()))

	children, err := repo.GetChildren(ctx, parentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d", len(children))
	}
}

func TestGetByPostID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := New(db)
	ctx := context.Background()

	postID := 7
	limit := 10
	offset := 0

	mock.ExpectQuery(`SELECT id, post_id, parent_id, content, created_at FROM comments WHERE post_id = \$1 AND parent_id IS NULL`).
		WithArgs(postID, limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "parent_id", "content", "created_at"}).
			AddRow(1, postID, nil, "top comment 1", time.Now()).
			AddRow(2, postID, nil, "top comment 2", time.Now()))

	comments, err := repo.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
}

func TestGetChildrenBatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := New(db)
	ctx := context.Background()

	parentIDs := []int{10, 20}

	mock.ExpectQuery(`SELECT id, post_id, parent_id, content, created_at FROM comments WHERE parent_id IN \(\$1, \$2\)`).
		WithArgs(parentIDs[0], parentIDs[1]).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "parent_id", "content", "created_at"}).
			AddRow(1, 100, 10, "child of 10", time.Now()).
			AddRow(2, 100, 10, "child of 10 again", time.Now()).
			AddRow(3, 100, 20, "child of 20", time.Now()))

	children, err := repo.GetChildrenBatch(ctx, parentIDs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(children) != 3 {
		t.Errorf("expected 3 children, got %d", len(children))
	}
}
