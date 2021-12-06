package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/jwt"
	postpb "github.com/os3224/final-project-0b5a2e16-babysuse/web/post/pb"
	userpb "github.com/os3224/final-project-0b5a2e16-babysuse/web/user/pb"
	"google.golang.org/grpc"
)

type Post struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	username := _Authenticate(&w, r)
	if len(username) == 0 {
		return
	}

	log.Printf("Getting posts")
	// get all followee (including self)
	followees := _GetFollowee(username)
	followees = append(followees, username)

	// set up RPC client
	conn, err := grpc.Dial(PostSrvAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := postpb.NewPostServiceClient(conn)

	// get posts from each followee (including self)
	var posts []Post
	for _, user := range followees {
		// contact RPC server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		resp, err := client.FollowPosts(ctx, &postpb.PostsRequest{Follower: user})
		if err != nil {
			log.Fatalf("failed to follow posts: %v", err)
		}
		for _, post := range resp.Posts {
			posts = append(posts, Post{
				Text:   post.Text,
				Author: post.Author,
			})
		}
	}
	json.NewEncoder(w).Encode(posts)
}

type FollowReq struct {
	Username    string `json:"username"`
	Unfollowing bool   `json:"unfollowing"`
}

func Follow(w http.ResponseWriter, r *http.Request) {
	username := _Authenticate(&w, r)
	if len(username) == 0 {
		return
	}

	// decode request body
	var following FollowReq
	err := json.NewDecoder(r.Body).Decode(&following)
	if err != nil {
		log.Printf("Invalid body format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// set up RPC client
	conn, err := grpc.Dial(UserSrvAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)

	// contact RPC server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if following.Unfollowing {
		_, err = client.Unfollow(ctx, &userpb.FollowRequest{Username: username, Followeename: following.Username})
	} else {
		_, err = client.Follow(ctx, &userpb.FollowRequest{Username: username, Followeename: following.Username})
	}
	if err != nil {
		log.Fatalf("failed to follow: %v", err)
	}

	// return the lastest following status
	followees := _GetFollowee(username)
	json.NewEncoder(w).Encode(followees)
}

func GetFollowees(w http.ResponseWriter, r *http.Request) {
	username := _Authenticate(&w, r)
	if len(username) == 0 {
		return
	}

	// get all followee
	log.Printf("Getting following")
	followees := _GetFollowee(username)
	json.NewEncoder(w).Encode(followees)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	username := _Authenticate(&w, r)
	if len(username) == 0 {
		return
	}

	// set up RPC client
	conn, err := grpc.Dial(UserSrvAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)

	// contact RPC server
	log.Printf("Getting users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rpcResp, err := client.GetUsers(ctx, &userpb.GetFolloweeRequest{Username: username})
	if err != nil {
		log.Fatalf("failed to get followee: %v", err)
	}

	// extract data
	var users []string
	for _, followee := range rpcResp.Followees {
		users = append(users, followee.Username)
	}
	json.NewEncoder(w).Encode(users)
}

func _Authenticate(w *http.ResponseWriter, r *http.Request) string {
	log.Printf("Authenticating")
	// authenticate user
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			(*w).WriteHeader(http.StatusUnauthorized)
			return ""
		}
		(*w).WriteHeader(http.StatusBadRequest)
		return ""
	}

	// validate jwt token
	username, err := jwt.ParseToken(cookie.Value)
	if err != nil {
		http.Error(*w, "Invalid cookie", http.StatusUnauthorized)
		return ""
	}

	return username
}

func _GetFollowee(username string) []string {
	// set up RPC client
	conn, err := grpc.Dial(UserSrvAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)

	// contact RPC server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetFollowee(ctx, &userpb.GetFolloweeRequest{Username: username})
	if err != nil {
		log.Fatalf("failed to get followee: %v", err)
	}

	// extract data
	var followees []string
	for _, followee := range resp.Followees {
		followees = append(followees, followee.Username)
	}
	return followees
}
