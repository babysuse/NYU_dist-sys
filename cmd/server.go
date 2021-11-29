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
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()
	database.Migrate()

	//mux := http.NewServeMux()
	middleware := auth.Middleware()
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	//mux.Handle("/query", middleware(server))
	//mux.Handle("/query", server)
	//http.Handle("/", playground.Handler("GraphQL server", "/query"))
	http.Handle("/query", middleware(server))

	log.Printf("connect to http://localhost:%s/ for GraphQL server", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
