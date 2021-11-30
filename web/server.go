package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"

	//"github.com/go-chi/chi"
	//"github.com/os3224/final-project-0b5a2e16-babysuse/internal/auth"

	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/auth"
	database "github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/db/migrations/mysql"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/account"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/generated"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/post"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/user"
)

const defaultPort = "16008"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// set up DB conn
	database.InitDB()
	database.Migrate()

	// set up RPC server
	account.NewAccountServiceServer()
	user.NewUserServiceServer()
	post.NewPostServiceServer()

	// set up web server
	middleware := auth.Middleware()
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/query", middleware(server))
	log.Printf("GraphQL server listening at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
