package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/yyyppd/graphql-comments-system/internal/comment"
)

type PostgresCommentRepo struct {
	DB *sql.DB
}

func New(db *sql.DB) *PostgresCommentRepo {
	return &PostgresCommentRepo{DB: db}
}

func (r *PostgresCommentRepo) Create(ctx context.Context, postID int, parentID *int, content string) (*comment.Comment, error) {
	var id int
	query := `
		INSERT INTO comments (post_id, parent_id, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var parent sql.NullInt64
	if parentID != nil {
		parent = sql.NullInt64{Int64: int64(*parentID), Valid: true}
	}

	err := r.DB.QueryRowContext(ctx, query, postID, parent, content).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert comment: %w", err)
	}

	return &comment.Comment{
		ID:        id,
		PostID:    postID,
		ParentID:  parent,
		Content:   content,
		CreatedAt: time.Now(), // можно заменить на возвращаемое значение из БД
	}, nil
}

func (r *PostgresCommentRepo) GetByPostID(ctx context.Context, postID, limit, offset int) ([]*comment.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, content, created_at
		FROM comments
		WHERE post_id = $1 AND parent_id IS NULL
		ORDER BY created_at
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select top-level comments: %w", err)
	}

	defer rows.Close()

	var result []*comment.Comment
	for rows.Next() {
		var c comment.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &c)
	}
	return result, nil
}

func (r *PostgresCommentRepo) GetChildren(ctx context.Context, parentID int) ([]*comment.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, content, created_at
		FROM comments
		WHERE parent_id = $1
		ORDER BY created_at
	`

	rows, err := r.DB.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("select child comments: %w", err)
	}

	defer rows.Close()

	var result []*comment.Comment
	for rows.Next() {
		var c comment.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &c)
	}
	return result, nil
}

func (r *PostgresCommentRepo) GetChildrenBatch(ctx context.Context, parentIDs []int) ([]*comment.Comment, error) {
	if len(parentIDs) == 0 {
		return []*comment.Comment{}, nil
	}

	// Подготавливаем параметры и плейсхолдеры: $1, $2, ...
	args := make([]interface{}, len(parentIDs))
	placeholders := make([]string, len(parentIDs))
	for i, id := range parentIDs {
		args[i] = id
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		SELECT id, post_id, parent_id, content, created_at
		FROM comments
		WHERE parent_id IN (%s)
		ORDER BY parent_id, created_at
	`,
		strings.Join(placeholders, ", "),
	)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select children batch: %w", err)
	}
	defer rows.Close()

	var comments []*comment.Comment
	for rows.Next() {
		var c comment.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}
