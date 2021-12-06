package main

import (
	"log"
	"net/http"
	"os"

	//"github.com/go-chi/chi"
	//"github.com/os3224/final-project-0b5a2e16-babysuse/internal/auth"

	database "github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/db/migrations/mysql"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/account"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/post"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/user"
	"github.com/rs/cors"
)

const (
	defaultPort    = "16008"
	AccountSrvAddr = "localhost:16018"
	UserSrvAddr    = "localhost:16028"
	PostSrvAddr    = "localhost:16038"
)

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
	// middleware := auth.Middleware()
	// server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	// http.Handle("/query", middleware(server))
	// log.Printf("GraphQL server listening at http://localhost:%s/", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/signup", Signup)
	mux.HandleFunc("/posts", GetPosts)
	mux.HandleFunc("/users", GetUsers)
	mux.HandleFunc("/following", GetFollowees)
	// mux.HandleFunc("/createpost", CreatePost)
	corsConfig := cors.New(cors.Options{
		AllowedHeaders:   []string{"Content-Type", "Cookies", "Origin"},
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		// Debug:            true,
	})
	corsHandler := corsConfig.Handler(mux)
	log.Printf("Web server listening at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
