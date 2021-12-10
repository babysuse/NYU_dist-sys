package post

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

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
	// stmt, err := database.DB.Prepare("INSERT INTO Posts(Text, Author) Values (?, ?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// insert data
	// result, err := stmt.Exec(post.Text, post.Author)
	// if err != nil {
	// 	return &pb.CreatePostResponse{}, err
	// }
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return &pb.CreatePostResponse{}, err
	// }
	// return &pb.CreatePostResponse{Post: &pb.Post{ID: strconv.FormatInt(id, 10), Text: post.Text, Author: post.Author}}, nil

	// Put to raft cluster
	client := http.Client{}
	data, err := json.Marshal(post.Text)
	if err != nil {
		log.Fatalf("Failed to Marshal: %v", err)
	}
	req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:16049/posts/"+post.Author, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Failed to NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		log.Fatalf("Failed to Do request: %v", err)
	}
	return &pb.CreatePostResponse{Post: &pb.Post{ID: "0", Text: post.Text, Author: post.Author}}, nil
}

func (srv *Server) FollowPosts(ctx context.Context, req *pb.PostsRequest) (*pb.PostsResponse, error) {
	// fetch data (via prepared SQL internally)
	// rows, err := database.DB.Query(`
	// 	SELECT ID, Text, Author
	// 	FROM Posts
	// 	WHERE Author = ?
	// `, req.Follower)
	// if err != nil {
	// 	return &pb.PostsResponse{}, err
	// }
	// defer rows.Close()

	// extract data
	// var resp pb.PostsResponse
	// for rows.Next() {
	// 	var post pb.Post
	// 	err := rows.Scan(&post.ID, &post.Text, &post.Author)
	// 	if err != nil {
	// 		return &pb.PostsResponse{}, err
	// 	}
	// 	resp.Posts = append(resp.Posts, &post)
	// }

	// get posts from raft cluster
	var resp pb.PostsResponse
	raftRespPosts, err := http.Get("http://127.0.0.1:16049/posts/" + req.Follower)
	if err != nil {
		log.Fatalf("Failed to Get: %v", err)
	}
	defer raftRespPosts.Body.Close()
	bytes, err := io.ReadAll(raftRespPosts.Body)
	if err != nil {
		log.Fatalf("Failed to ReadAll: %v", err)
	}
	var posts []string
	json.Unmarshal(bytes, &posts)
	log.Printf("%s's posts: %v", req.Follower, posts)
	// prepare response
	for _, text := range posts {
		resp.Posts = append(resp.Posts, &pb.Post{Text: text, Author: req.Follower})
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
