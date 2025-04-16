package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yyyppd/graphql-comments-system/internal/post"
)

type PostgresPostRepo struct {
	DB *sql.DB
}

func New(db *sql.DB) *PostgresPostRepo {
	return &PostgresPostRepo{DB: db}
}

func (r *PostgresPostRepo) GetAll(ctx context.Context) ([]*post.Post, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, title, content, allow_comments FROM posts`)
	if err != nil {
		return nil, fmt.Errorf("query posts: %w", err)
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		var p post.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AllowComments); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

func (r *PostgresPostRepo) GetByID(ctx context.Context, id int) (*post.Post, error) {
	var p post.Post
	err := r.DB.QueryRowContext(ctx, `SELECT id, title, content, allow_comments FROM posts WHERE id = $1`, id).
		Scan(&p.ID, &p.Title, &p.Content, &p.AllowComments)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get post by id: %w", err)
	}
	return &p, nil
}

func (r *PostgresPostRepo) Create(ctx context.Context, title, content string, allowComments bool) (*post.Post, error) {
	var id int
	err := r.DB.QueryRowContext(ctx, `
		INSERT INTO posts (title, content, allow_comments)
		VALUES ($1, $2, $3)
		RETURNING id`, title, content, allowComments).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert post: %w", err)
	}
	return &post.Post{
		ID:            id,
		Title:         title,
		Content:       content,
		AllowComments: allowComments,
	}, nil
}
