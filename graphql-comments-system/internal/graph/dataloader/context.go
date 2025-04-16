package dataloader

import (
	"context"
	"net/http"
)

type ctxKey string

const loadersKey = ctxKey("dataloaders")

func Middleware(loaders *Loaders) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func For(ctx context.Context) *Loaders {
	val := ctx.Value(loadersKey)
	if val == nil {
		panic("No dataloader in context!")
	}
	return val.(*Loaders)
}
