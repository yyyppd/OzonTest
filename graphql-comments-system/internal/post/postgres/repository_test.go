package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreatePost(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := New(db)
	ctx := context.Background()

	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs("New Title", "New Content", true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	post, err := repo.Create(ctx, "New Title", "New Content", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.ID != 10 || post.Title != "New Title" {
		t.Errorf("post not created properly: %+v", post)
	}
}

func TestGetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := New(db)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT id, title, content, allow_comments FROM posts WHERE id = \$1`).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "allow_comments"}).
			AddRow(7, "Test", "Content", true))

	p, err := repo.GetByID(ctx, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p == nil || p.ID != 7 || !p.AllowComments {
		t.Errorf("unexpected post result: %+v", p)
	}
}

func TestGetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := New(db)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT id, title, content, allow_comments FROM posts`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "allow_comments"}).
			AddRow(1, "A", "X", true).
			AddRow(2, "B", "Y", false))

	all, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(all) != 2 {
		t.Errorf("expected 2 posts, got %d", len(all))
	}
}
