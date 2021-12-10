package account

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

	database "github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/db/migrations/mysql"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/jwt"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/account/autherr"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/account/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

const (
	port = ":16018"
)

type Server struct {
	pb.UnimplementedAccountServiceServer
}

func (srv *Server) Signup(ctx context.Context, acc *pb.Account) (*pb.Token, error) {
	// prepare SQL
	statement, err := database.DB.Prepare("INSERT INTO Users(Username, Password) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	// insert data
	hashedPassword, err := HashPassword(acc.Password)
	_, err = statement.Exec(acc.Username, hashedPassword)
	if err != nil {
		log.Println(err)
		return &pb.Token{}, err
	}

	// generate token
	token, err := jwt.GenerateToken(acc.Username)
	if err != nil {
		return &pb.Token{}, err
	}
	log.Printf("New user: %s\n", acc.Username)
	return &pb.Token{Token: token}, nil
}

func (srv *Server) Login(ctx context.Context, acc *pb.Account) (*pb.Token, error) {
	// prepare SQL
	// if err := database.DB.Ping(); err != nil {
	// 	log.Panic(err)
	// }
	// stmt, err := database.DB.Prepare("SELECT Password FROM Users WHERE Username = ?")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fetch data
	// row := stmt.QueryRow(acc.Username)
	// var hashedPassword string
	// err = row.Scan(&hashedPassword)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return &pb.Token{}, &autherr.WrongAuth{}
	// 	} else {
	// 		log.Fatal(err)
	// 	}
	// }

	// GET from raft clusters
	resp, err := http.Get("http://127.0.0.1:16049/users/" + acc.Username)
	if err != nil {
		log.Fatalf("Failed to Get: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to ReadAll: %v", err)
	}
	var hashedPassword string
	json.Unmarshal(bytes, &hashedPassword)

	// authenticate
	if !CheckPasswordHash(acc.Password, hashedPassword) {
		return &pb.Token{}, &autherr.WrongAuth{}
	}

	// generate token
	token, err := jwt.GenerateToken(acc.Username)
	if err != nil {
		return &pb.Token{}, err
	}
	return &pb.Token{Token: token}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.DB.Prepare("SELECT ID FROM Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var id int
	err = row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return 0, err
	}
	return id, nil
}

func NewAccountServiceServer() {
	go func() {
		bind, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		srv := grpc.NewServer()
		pb.RegisterAccountServiceServer(srv, &Server{})
		log.Printf("account service listening at %v", bind.Addr())
		if err := srv.Serve(bind); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
