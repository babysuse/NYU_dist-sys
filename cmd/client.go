package main

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

var client = graphql.NewClient("http://localhost:8080/query")
var token string

//func TestCreateUser(t *testing.T) {
func TestCreateUser() string {
	request := graphql.NewRequest(`
		mutation {
			signup(input: {username: "user", password: "passwd"})
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
	client := graphql.NewClient("http://localhost:8080/query")
	request := graphql.NewRequest(`
		mutation {
			login(input: {username: "user", password: "passwd"})
		}
	`)

	var response struct {
		Login string `json:"login"`
	}
	if err := client.Run(context.Background(), request, &response); err != nil {
		panic(err)
	}
	return response.Login
}

func TestFollow() string {
	request := graphql.NewRequest(`
		mutation {
			follow(input: {userid: "12"})
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

func TestUnfollow() string {
	request := graphql.NewRequest(`
		mutation {
			unfollow(input: {userid: "12"})
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

//func TestCreatePost(t *testing.T) {
//func TestCreatePost(t *testing.T) {
func TestCreatePost() {
	request := graphql.NewRequest(`
		mutation {
			createPost(input: {text: "new content created!!"}) {
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
func TestPosts() {
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
	//got := response.Posts[0].Author.Name
	//want := "babysuse"
	//if got != want {
	//t.Errorf("Expect %q, got %q", want, got)
	//}
}

func main() {
	//token = TestCreateUser()
	token = UserLogin()
	println(token)
	//TestFollow()
	TestUnfollow()
	//TestCreatePost()
	TestPosts()
}
