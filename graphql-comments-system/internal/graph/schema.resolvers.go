package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"context"
	"fmt"
	"strconv"
	"time"

	godataloader "github.com/graph-gophers/dataloader"
	"github.com/yyyppd/graphql-comments-system/internal/comment"
	localdataloader "github.com/yyyppd/graphql-comments-system/internal/graph/dataloader"
	"github.com/yyyppd/graphql-comments-system/internal/graph/model"
	"github.com/yyyppd/graphql-comments-system/pkg/utils"
)

// Children is the resolver for the children field.
func (r *commentResolver) Children(ctx context.Context, obj *model.Comment) ([]*model.Comment, error) {
	loaders := localdataloader.For(ctx)

	thunk := loaders.ChildrenByParentID.Load(ctx, godataloader.StringKey(obj.ID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	rawComments, ok := result.([]*comment.Comment)
	if !ok {
		return nil, fmt.Errorf("invalid data from dataloader")
	}

	var gqlComments []*model.Comment
	for _, c := range rawComments {
		gqlComments = append(gqlComments, &model.Comment{
			ID:        strconv.Itoa(c.ID),
			PostID:    strconv.Itoa(c.PostID),
			ParentID:  utils.PtrToString(c.ParentID),
			Content:   c.Content,
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
			Children:  []*model.Comment{},
		})
	}

	return gqlComments, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, allowComments bool) (*model.Post, error) {
	p, err := r.PostRepo.Create(ctx, title, content, allowComments)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:            strconv.Itoa(p.ID),
		Title:         p.Title,
		Content:       p.Content,
		AllowComments: p.AllowComments,
		Comments:      []*model.Comment{},
	}, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*model.Comment, error) {
	pid, err := strconv.Atoi(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID")
	}

	post, err := r.PostRepo.GetByID(ctx, pid)
	if err != nil || post == nil {
		return nil, fmt.Errorf("post not found")
	}
	if !post.AllowComments {
		return nil, fmt.Errorf("comments are disabled for this post")
	}

	var parentInt *int
	if parentID != nil {
		id, err := strconv.Atoi(*parentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent ID")
		}
		parentInt = &id
	}

	c, err := r.Resolver.CommentRepo.Create(ctx, pid, parentInt, content)
	if err != nil {
		return nil, err
	}

	gqlComment := &model.Comment{
		ID:        strconv.Itoa(c.ID),
		PostID:    strconv.Itoa(c.PostID),
		ParentID:  utils.PtrToString(c.ParentID),
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		Children:  []*model.Comment{},
	}

	r.Resolver.mu.Lock()
	subs := r.Resolver.commentSubscribers[strconv.Itoa(pid)]
	for _, ch := range subs {
		select {
		case ch <- gqlComment:
		default:
			// если канал занят — пропускаем, чтобы не блокировать
		}
	}
	r.Resolver.mu.Unlock()

	return gqlComment, nil
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *model.Post, limit *int32, offset *int32) ([]*model.Comment, error) {
	postID, err := strconv.Atoi(obj.ID)
	if err != nil {
		return nil, err
	}

	l := 10
	if limit != nil {
		l = int(*limit)
	}
	o := 0
	if offset != nil {
		o = int(*offset)
	}

	comments, err := r.Resolver.CommentRepo.GetByPostID(ctx, postID, l, o)
	if err != nil {
		return nil, err
	}

	var result []*model.Comment
	for _, c := range comments {
		result = append(result, &model.Comment{
			ID:        strconv.Itoa(c.ID),
			PostID:    strconv.Itoa(c.PostID),
			ParentID:  utils.PtrToString(c.ParentID),
			Content:   c.Content,
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
			Children:  []*model.Comment{}, // рекурсия?
		})
	}

	return result, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.PostRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var gqlPosts []*model.Post
	for _, p := range posts {
		gqlPosts = append(gqlPosts, &model.Post{
			ID:            strconv.Itoa(p.ID),
			Title:         p.Title,
			Content:       p.Content,
			AllowComments: p.AllowComments,
			Comments:      []*model.Comment{},
		})
	}
	return gqlPosts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	p, err := r.PostRepo.GetByID(ctx, postID)
	if err != nil || p == nil {
		return nil, err
	}

	return &model.Post{
		ID:            strconv.Itoa(p.ID),
		Title:         p.Title,
		Content:       p.Content,
		AllowComments: p.AllowComments,
		Comments:      []*model.Comment{},
	}, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	commentChan := make(chan *model.Comment, 1)

	r.mu.Lock()
	r.commentSubscribers[postID] = append(r.commentSubscribers[postID], commentChan)
	r.mu.Unlock()

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		subs := r.commentSubscribers[postID]
		for i, ch := range subs {
			if ch == commentChan {
				r.commentSubscribers[postID] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		r.mu.Unlock()
	}()

	return commentChan, nil
}

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
