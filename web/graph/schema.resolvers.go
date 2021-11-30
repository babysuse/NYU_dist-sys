package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/auth"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/jwt"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/posts"
	autherr "github.com/os3224/final-project-0b5a2e16-babysuse/web/account/autherr"
	accountpb "github.com/os3224/final-project-0b5a2e16-babysuse/web/account/pb"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/generated"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/model"
	userpb "github.com/os3224/final-project-0b5a2e16-babysuse/web/user/pb"
	"google.golang.org/grpc"
)

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

func (r *mutationResolver) Signup(ctx context.Context, input *model.Signup) (string, error) {
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
	resp, err := client.Signup(ctx, &accountpb.Account{Username: input.Username, Password: input.Password})
	if err != nil {
		log.Fatalf("failed to signup: %v", err)
	}
	return resp.Token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input *model.Login) (string, error) {
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
	resp, err := client.Login(ctx, &accountpb.Account{Username: input.Username, Password: input.Password})
	if err != nil {
		_, wrongtype := err.(*autherr.WrongAuth)
		if !wrongtype {
			return "", err
		}
		log.Fatalf("failed to login: %v", err)
	}
	return resp.Token, nil
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

func (r *mutationResolver) Follow(ctx context.Context, input *model.Followee) (string, error) {
	// authenticate user
	user := auth.ForContext(ctx)
	if user == nil {
		return "", fmt.Errorf("Access denied")
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
	_, err = client.Follow(ctx, &userpb.FollowRequest{Username: user.Username, Followeename: input.Username})
	if err != nil {
		log.Fatalf("failed to follow: %v", err)
	}
	return "Followed", nil
}

func (r *mutationResolver) Unfollow(ctx context.Context, input *model.Followee) (string, error) {
	// authenticate user
	user := auth.ForContext(ctx)
	if user == nil {
		return "", fmt.Errorf("Access denied")
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
	_, err = client.Unfollow(ctx, &userpb.FollowRequest{Username: user.Username, Followeename: input.Username})
	if err != nil {
		log.Fatalf("failed to unfollow: %v", err)
	}
	return "Unfollowed", nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	user := auth.ForContext(ctx)
	var gqlPosts []*model.Post
	if user == nil {
		return gqlPosts, fmt.Errorf("Access denied")
	}
	// get all followee (including self)
	//followees := user.GetFollowee()
	followees := []accountpb.Account{}
	followees = append(followees, *user)

	// get posts from each followee (including self)
	var dbPosts []posts.Post
	for _, user := range followees {
		dbPosts = posts.GetAll(user.Username)
		for _, post := range dbPosts {
			gqlUser := &model.User{ID: post.User.ID, Name: post.User.Username}
			gqlPosts = append(gqlPosts, &model.Post{
				ID:     post.ID,
				Text:   post.Text,
				Author: gqlUser,
			})
		}
	}
	return gqlPosts, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
const (
	AccountSrvAddr = "localhost:16018"
	UserSrvAddr    = "localhost:16028"
)
