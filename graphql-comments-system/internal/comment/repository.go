package comment

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	ParentID  sql.NullInt64
	Content   string
	CreatedAt time.Time
}

type CommentRepository interface {
	Create(ctx context.Context, postID int, parentID *int, content string) (*Comment, error)
	GetByPostID(ctx context.Context, postID, limit, offset int) ([]*Comment, error)
	GetChildren(ctx context.Context, parentID int) ([]*Comment, error)
	GetChildrenBatch(ctx context.Context, parentIDs []int) ([]*Comment, error)
}
