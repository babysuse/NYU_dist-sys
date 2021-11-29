package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/auth"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/jwt"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/posts"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/users"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/generated"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/model"
)

func (r *mutationResolver) Signup(ctx context.Context, input *model.Signup) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	fmt.Printf("New user: %s\n", user.Username)
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input *model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		return "", &users.WrongAuth{}
	}
	fmt.Printf("%s logged in\n", user.Username)

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshToken) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("Access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, err
}

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Post{}, fmt.Errorf("Access denied")
	}
	var post posts.Post
	post.User = user
	post.Text = input.Text
	postID := post.Save()
	return &model.Post{ID: strconv.FormatInt(postID, 10), Text: post.Text}, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	user := auth.ForContext(ctx)
	var gqlPosts []*model.Post
	if user == nil {
		return gqlPosts, fmt.Errorf("Access denied")
	}

	var dbPosts []posts.Post
	dbPosts = posts.GetAll(user.ID)
	for _, post := range dbPosts {
		gqlUser := &model.User{ID: post.User.ID, Name: post.User.Username}
		gqlPosts = append(gqlPosts, &model.Post{
			ID:     post.ID,
			Text:   post.Text,
			Author: gqlUser,
		})
	}
	return gqlPosts, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
