package main

import (
	"log"
	"os/exec"
	"strconv"

	"net/http"
	"os"

	//"github.com/os3224/final-project-0b5a2e16-babysuse/internal/auth"
	// database "github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/db/migrations/mysql"
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

var _rafrSrvPort = []string{"16048", "16058", "16068"}
var raftSrvPort = []string{"16049", "16059", "16069"}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// set up DB conn
	// database.InitDB()
	// database.Migrate()

	cluster := ""
	for _, port := range _rafrSrvPort {
		cluster += "http://127.0.0.1:" + port + ","
	}
	cluster = cluster[:len(cluster)-1]

	// set up raft cluster
	//*
	for i, port := range raftSrvPort {
		go func(id, cluster, port string) {
			// ./raft --id 1 --cluster http://127.0.0.1:12345 --port 12346
			err := exec.Command("./raft", "--id", id, "--cluster", cluster, "--port", port).Run()
			if err != nil {
				log.Fatalf("fail to launch raft node %s: %s", id, err)
			}
		}(strconv.Itoa(i+1), cluster, port)
	}
	// */

	// set up RPC server
	account.NewAccountServiceServer()
	user.NewUserServiceServer()
	post.NewPostServiceServer()

	// set up graphQL web server
	// middleware := auth.Middleware()
	// server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	// http.Handle("/query", middleware(server))
	// log.Printf("GraphQL server listening at http://localhost:%s/", port)

	// set up REST web server
	mux := http.NewServeMux()
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/signup", Signup)
	mux.HandleFunc("/get_posts", GetPosts)
	mux.HandleFunc("/get_users", GetUsers)
	mux.HandleFunc("/get_following", GetFollowees)
	mux.HandleFunc("/follow", Follow)
	mux.HandleFunc("/createpost", CreatePost)
	corsConfig := cors.New(cors.Options{
		AllowedHeaders:   []string{"Content-Type", "Cookies", "Origin"},
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002"},
		AllowCredentials: true,
		// Debug:            true,
	})
	corsHandler := corsConfig.Handler(mux)
	log.Printf("Web server listening at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
