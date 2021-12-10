package user

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

	// database "github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/db/migrations/mysql"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/user/pb"
	"google.golang.org/grpc"
)

const (
	port = ":16028"
)

type Server struct {
	pb.UnimplementedUserServiceServer
}

func (srv *Server) Follow(ctx context.Context, req *pb.FollowRequest) (*pb.FollowReply, error) {
	// insert database
	// stmt, err := database.DB.Prepare("INSERT INTO Following(Username, Followeename) VALUES(?, ?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = stmt.Exec(req.Username, req.Followeename)
	// if err != nil {
	// 	log.Println(err)
	// 	return &pb.FollowReply{}, err
	// }

	// Put to raft cluster
	client := http.Client{}
	data, err := json.Marshal(req.Followeename)
	if err != nil {
		log.Fatalf("Failed to Marshal: %v", err)
	}
	raftReq, err := http.NewRequest(
		http.MethodPut,
		"http://127.0.0.1:16049/following/"+req.Username,
		bytes.NewBuffer(data),
	)
	if err != nil {
		log.Fatalf("Failed to NewRequest: %v", err)
	}
	_, err = client.Do(raftReq)
	if err != nil {
		log.Fatalf("Failed to Do request: %v", err)
	}
	return &pb.FollowReply{Result: "Followed"}, nil
}

func (srv *Server) Unfollow(ctx context.Context, req *pb.FollowRequest) (*pb.FollowReply, error) {
	// stmt, err := database.DB.Prepare("DELETE FROM Following WHERE Username = ? AND Followeename = ?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = stmt.Exec(req.Username, req.Followeename)
	// if err != nil {
	// 	log.Println(err)
	// 	return &pb.FollowReply{}, err
	// }

	// Put to raft cluster
	client := http.Client{}
	data, err := json.Marshal(req.Followeename)
	if err != nil {
		log.Fatalf("Failed to Marshal: %v", err)
	}
	raftReq, err := http.NewRequest(
		http.MethodPut,
		"http://127.0.0.1:16049/following/un/"+req.Username,
		bytes.NewBuffer(data),
	)
	if err != nil {
		log.Fatalf("Failed to NewRequest: %v", err)
	}
	_, err = client.Do(raftReq)
	if err != nil {
		log.Fatalf("Failed to Do request: %v", err)
	}
	return &pb.FollowReply{Result: "Unfollowed"}, nil
}

func (srv *Server) GetFollowee(ctx context.Context, req *pb.GetFolloweeRequest) (*pb.GetFolloweeResponse, error) {
	// prepare SQL
	// rows, err := database.DB.Query(`
	// 	SELECT Followeename
	// 	FROM Following
	// 	WHERE Username = ?`, req.Username)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// fetch data
	// var resp pb.GetFolloweeResponse
	// for rows.Next() {
	// 	var followee pb.User
	// 	err := rows.Scan(&followee.Username)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	resp.Followees = append(resp.Followees, &followee)
	// }

	// get from raft cluster
	raftResp, err := http.Get("http://127.0.0.1:16049/following/" + req.Username)
	if err != nil {
		log.Fatalf("Failed to Get: %v", err)
	}
	defer raftResp.Body.Close()
	bytes, err := io.ReadAll(raftResp.Body)
	if err != nil {
		log.Fatalf("Failed to ReadAll: %v", err)
	}
	var users []string
	json.Unmarshal(bytes, &users)

	// prepare response
	var resp pb.GetFolloweeResponse
	for _, u := range users {
		followee := pb.User{Username: u}
		resp.Followees = append(resp.Followees, &followee)
	}
	return &resp, nil
}

func (srv *Server) GetUsers(ctx context.Context, in *pb.GetFolloweeRequest) (*pb.GetFolloweeResponse, error) {
	// query database
	// rows, err := database.DB.Query(`
	// 	SELECT Username
	// 	FROM Users
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// extract data
	// var resp pb.GetFolloweeResponse
	// for rows.Next() {
	// 	var user pb.User
	// 	err := rows.Scan(&user.Username)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	resp.Followees = append(resp.Followees, &user)
	// }

	// get from raft cluster
	raftResp, err := http.Get("http://127.0.0.1:16049/users")
	if err != nil {
		log.Fatalf("Failed to Get: %v", err)
	}
	defer raftResp.Body.Close()
	bytes, err := io.ReadAll(raftResp.Body)
	if err != nil {
		log.Fatalf("Failed to ReadAll: %v", err)
	}
	var users []string
	json.Unmarshal(bytes, &users)

	// prepare response
	var resp pb.GetFolloweeResponse
	for _, u := range users {
		user := pb.User{Username: u}
		resp.Followees = append(resp.Followees, &user)
	}
	return &resp, nil
}

func NewUserServiceServer() {
	go func() {
		bind, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to bind: %v", err)
		}
		srv := grpc.NewServer()
		pb.RegisterUserServiceServer(srv, &Server{})
		log.Printf("user service listening at %v", bind.Addr())
		if err := srv.Serve(bind); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
