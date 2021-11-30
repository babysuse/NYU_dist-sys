package main

import (
	"context"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
)

var client = graphql.NewClient("http://localhost:16008/query")
var token string

// user:passwd
// test:test123
// follower:follower
func CreateUser() string {
	request := graphql.NewRequest(`
		mutation {
			signup(input: {username: "follower", password: "follower"})
		}
	`)

	var response struct {
		Signup string `json:"signup"`
	}
	if err := client.Run(context.Background(), request, &response); err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", response)
	return response.Signup
}

func UserLogin() string {
	//login(input: {username: "user", password: "passwd"})
	request := graphql.NewRequest(`
		mutation {
			login(input: {username: "follower", password: "follower"})
		}
	`)

	var response struct {
		Login string `json:"login"`
	}
	if err := client.Run(context.Background(), request, &response); err != nil {
		log.Fatalf("fail to login: %v", err)
	}
	return response.Login
}

func Follow() string {
	request := graphql.NewRequest(`
		mutation {
			follow(input: {username: "test"})
		}
	`)
	request.Header.Set("Authorization", token)

	var response struct {
		Follow string `json:"follow"`
	}
	if err := client.Run(context.Background(), request, &response); err != nil {
		panic(err)
	}
	return response.Follow
}

func Unfollow() string {
	request := graphql.NewRequest(`
		mutation {
			unfollow(input: {username: "user"})
		}
	`)
	request.Header.Set("Authorization", token)

	var response struct {
		Unfollow string `json:"unfollow"`
	}
	if err := client.Run(context.Background(), request, &response); err != nil {
		panic(err)
	}
	return response.Unfollow
}

func CreatePost() {
	request := graphql.NewRequest(`
		mutation {
			createPost(input: {text: "Created by user"}) {
				id,
				text,
			}
		}
	`)
	request.Header.Set("Authorization", token)

	var response struct {
		CreatePost struct {
			ID   string `json:"ID"`
			Text string `json:"Text"`
		} `json:"createPost"`
	}
	if err := client.Run(context.Background(), request, &response); err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", response)
	//want := "new content created!!"
	//got := response.CreatePost.Text
	//if got != want {
	//t.Errorf("Expect %q, got %q", want, got)
}

//func TestPosts(t *testing.T) {
func Posts() {
	request := graphql.NewRequest(`
		query {
			posts {
				id,
				text,
				author {
					name
				}
			}
		}
	`)
	request.Header.Set("Authorization", token)

	var response struct {
		Posts []struct {
			Author struct {
				Name string `json:"name"`
			} `json:"author"`
			Text string `json:"text"`
			ID   string `json:"id"`
		} `json:"posts"`
	}
	if err := client.Run(context.Background(), request, &response); err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", response)
}

func main() {
	//token = CreateUser()
	token = UserLogin()
	println(token)
	//Unfollow()
	//Follow()
	//CreatePost()
	Posts()
}
