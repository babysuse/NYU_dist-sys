package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/auth"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/jwt"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/account/autherr"
	accountpb "github.com/os3224/final-project-0b5a2e16-babysuse/web/account/pb"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/generated"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/graph/model"
	postpb "github.com/os3224/final-project-0b5a2e16-babysuse/web/post/pb"
	userpb "github.com/os3224/final-project-0b5a2e16-babysuse/web/user/pb"
	"google.golang.org/grpc"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	// authenticate user
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Post{}, fmt.Errorf("Access denied")
	}

	// set up RPC client
	conn, err := grpc.Dial(PostSrvAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := postpb.NewPostServiceClient(conn)

	// contact RPC server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := client.CreatePost(ctx, &postpb.CreatePostRequest{Text: input.Text, Author: user.Username})
	if err != nil {
		log.Fatalf("failed to create post: %v", err)
	}
	return &model.Post{ID: resp.Post.ID, Text: resp.Post.Text}, nil
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

func (r *mutationResolver) Follow(ctx context.Context, input *model.FolloweeInput) (string, error) {
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

func (r *mutationResolver) Unfollow(ctx context.Context, input *model.FolloweeInput) (string, error) {
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

func (r *queryResolver) Getfollowee(ctx context.Context) ([]*model.Followee, error) {
	// authenticate user
	user := auth.ForContext(ctx)
	if user == nil {
		return []*model.Followee{}, fmt.Errorf("Access denied")
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
	resp, err := client.GetFollowee(ctx, &userpb.GetFolloweeRequest{Username: user.Username})
	if err != nil {
		log.Fatalf("failed to get followee: %v", err)
	}

	// extract data
	var followees []*model.Followee
	for _, followee := range resp.Followees {
		followees = append(followees, &model.Followee{Username: followee.Username})
	}
	return followees, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	// authenticate user
	user := auth.ForContext(ctx)
	if user == nil {
		return []*model.Post{}, fmt.Errorf("Access denied")
	}

	// get all followee (including self)
	followees, err := r.Getfollowee(ctx)
	followees = append(followees, &model.Followee{Username: user.Username})
	if err != nil {
		log.Fatalf("failed to get followee: %v", err)
	}

	// set up RPC client
	conn, err := grpc.Dial(PostSrvAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := postpb.NewPostServiceClient(conn)

	// get posts from each followee (including self)
	var gqlPosts []*model.Post
	for _, user := range followees {
		// contact RPC server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		resp, err := client.FollowPosts(ctx, &postpb.PostsRequest{Follower: user.Username})
		if err != nil {
			log.Fatalf("failed to follow posts: %v", err)
		}
		for _, post := range resp.Posts {
			gqlPosts = append(gqlPosts, &model.Post{
				ID:     post.ID,
				Text:   post.Text,
				Author: &model.User{Name: post.Author},
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
	PostSrvAddr    = "localhost:16038"
)
