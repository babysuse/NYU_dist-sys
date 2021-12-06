package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	accountpb "github.com/os3224/final-project-0b5a2e16-babysuse/web/account/pb"
	"google.golang.org/grpc"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	// decode request
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("%s logging in", cred.Username)

	// set up RPC client
	conn, err := grpc.Dial(AccountSrvAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := accountpb.NewAccountServiceClient(conn)

	// contact RPC server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := client.Signup(ctx, &accountpb.Account{Username: cred.Username, Password: cred.Password})
	if err != nil {
		log.Fatalf("failed to signup: %v", err)
	}

	// set cookie with expiry time of 10 minutes
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   resp.Token,
		Expires: time.Now().Add(10 * time.Minute),
	})
	log.Printf("%s signed up", cred.Username)
}
