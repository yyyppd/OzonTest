package dataloader

import (
	"context"
	"strconv"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/yyyppd/graphql-comments-system/internal/comment"
)

// This loader batch-loads children for a slice of parent comment IDs
type Loaders struct {
	ChildrenByParentID *dataloader.Loader
}

type CommentReader interface {
	GetChildren(ctx context.Context, parentID int) ([]*comment.Comment, error)
	GetChildrenBatch(ctx context.Context, parentIDs []int) ([]*comment.Comment, error)
}

func NewLoaders(repo CommentReader) *Loaders {
	return &Loaders{
		ChildrenByParentID: dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
			var parentIDs []int
			keyOrder := make(map[int]int)
			for i, key := range keys {
				id, err := strconv.Atoi(key.String())
				if err != nil {
					return []*dataloader.Result{{Error: err}}
				}
				parentIDs = append(parentIDs, id)
				keyOrder[id] = i
			}

			allChildren, err := repo.GetChildrenBatch(ctx, parentIDs)
			if err != nil {
				results := make([]*dataloader.Result, len(keys))
				for i := range keys {
					results[i] = &dataloader.Result{Error: err}
				}
				return results
			}

			grouped := make(map[int][]*comment.Comment)
			for _, c := range allChildren {
				parent := int(c.ParentID.Int64)
				grouped[parent] = append(grouped[parent], c)
			}

			results := make([]*dataloader.Result, len(keys))
			for i, key := range keys {
				id, _ := strconv.Atoi(key.String())
				results[i] = &dataloader.Result{Data: grouped[id]}
			}

			return results
		},
			dataloader.WithWait(5*time.Millisecond),
			dataloader.WithBatchCapacity(100),
		),
	}
}
