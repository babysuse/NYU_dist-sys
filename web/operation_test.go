package main

import (
	"testing"
	"time"
)

/*
 * Default state (USER : # OF POST : FOLLOWINGS):
 * 		test : 1 : {}
 * 		user : 1 : {}
 * 		follower : 0 : {test user}
 */

func TestGetPosts(t *testing.T) {
	posts := _GetPosts("test")
	if len(posts) != 1 {
		t.Errorf("Expect 1 post, got %d", len(posts))
	}
}

func TestCreatePost(t *testing.T) {
	_CreatePost("test", "2nd post by test")

	// wait for raft cluster to commit the update
	time.Sleep(300 * time.Millisecond)
	posts := _GetPosts("test")
	if len(posts) != 2 {
		t.Errorf("Expect 2 posts, got %d", len(posts))
	}
}

func TestGetFollowees(t *testing.T) {
	following := _GetFollowee("follower")
	if len(following) != 2 {
		t.Errorf("Expect 2 following, got %d", len(following))
	}
}

func TestGetUsers(t *testing.T) {
	users := _GetUsers("user")
	// three default user plus a user created by TestSignup()
	if len(users) != 4 {
		t.Errorf("Expect 4 users, got %d", len(users))
	}
}

func TestFollow(t *testing.T) {
	_Follow("follower", "user", true)

	// wait for raft cluster to commit the update
	time.Sleep(300 * time.Millisecond)
	following := _GetFollowee("follower")
	if len(following) != 1 {
		t.Errorf("Expect 1 following, got %d", len(following))
	}
}
