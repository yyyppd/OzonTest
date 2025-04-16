package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/yyyppd/graphql-comments-system/internal/comment"
	commentMemory "github.com/yyyppd/graphql-comments-system/internal/comment/memory"
	commentPostgres "github.com/yyyppd/graphql-comments-system/internal/comment/postgres"
	"github.com/yyyppd/graphql-comments-system/internal/graph"

	"github.com/yyyppd/graphql-comments-system/internal/post"
	postMemory "github.com/yyyppd/graphql-comments-system/internal/post/memory"
	postPostgres "github.com/yyyppd/graphql-comments-system/internal/post/postgres"

	"github.com/yyyppd/graphql-comments-system/internal/db"
	"github.com/yyyppd/graphql-comments-system/internal/graph/dataloader"
)

const defaultPort = "8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using defaults/environment")
	}

	useMemory := os.Getenv("USE_IN_MEMORY") == "true"

	var postRepo post.PostRepository
	var commentRepo comment.CommentRepository

	if useMemory {
		log.Println("Using in-memory repositories")
		postRepo = postMemory.New()
		commentRepo = commentMemory.New()
	} else {
		log.Println("Using PostgreSQL repositories")
		db.InitPostgres()
		postRepo = postPostgres.New(db.DB)
		commentRepo = commentPostgres.New(db.DB)
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(postRepo, commentRepo),
	}))

	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	loaders := dataloader.NewLoaders(commentRepo)

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", dataloader.Middleware(loaders)(srv))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Server ready at http://localhost:%s/ (in-memory: %v)", port, useMemory)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
