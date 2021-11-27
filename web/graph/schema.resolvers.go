package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/jwt"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/posts"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/users"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/generated"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/model"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	var post posts.Post
	post.Text = input.Text
	post.User.ID = input.AuthorID
	postID := post.Save()
	return &model.Post{ID: strconv.FormatInt(postID, 10), Text: post.Text}, nil
}

func (r *mutationResolver) Signup(ctx context.Context, input *model.Signup) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
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

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	var posts []*model.Post
	dummyPost := model.Post{
		ID:     "postid",
		Text:   "content",
		Author: &model.User{Name: "babysuse"},
	}
	posts = append(posts, &dummyPost)
	return posts, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
