package post

import (
	"context"
	"log"
	"net"
	"strconv"

	database "github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/db/migrations/mysql"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/post/pb"
	"google.golang.org/grpc"
)

const (
	port = ":16038"
)

type Server struct {
	pb.UnimplementedPostServiceServer
}

func (srv *Server) CreatePost(ctx context.Context, post *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	// prepare SQL
	stmt, err := database.DB.Prepare("INSERT INTO Posts(Text, Author) Values (?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	// insert data
	result, err := stmt.Exec(post.Text, post.Author)
	if err != nil {
		return &pb.CreatePostResponse{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return &pb.CreatePostResponse{}, err
	}
	return &pb.CreatePostResponse{Post: &pb.Post{ID: strconv.FormatInt(id, 10), Text: post.Text, Author: post.Author}}, nil
}

func (srv *Server) FollowPosts(ctx context.Context, req *pb.PostsRequest) (*pb.PostsResponse, error) {
	// fetch data (via prepared SQL internally)
	rows, err := database.DB.Query(`
		SELECT ID, Text, Author
		FROM Posts
		WHERE Author = ?
	`, req.Follower)
	if err != nil {
		return &pb.PostsResponse{}, err
	}
	defer rows.Close()

	// extract data
	var resp pb.PostsResponse
	for rows.Next() {
		var post pb.Post
		err := rows.Scan(&post.ID, &post.Text, &post.Author)
		if err != nil {
			return &pb.PostsResponse{}, err
		}
		resp.Posts = append(resp.Posts, &post)
	}
	return &resp, nil
}

func NewPostServiceServer() {
	go func() {
		bind, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to bind: %v", err)
		}
		srv := grpc.NewServer()
		pb.RegisterPostServiceServer(srv, &Server{})
		log.Printf("post service listening at %v", bind.Addr())
		if err := srv.Serve(bind); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
