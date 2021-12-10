package main

import (
	"testing"
	"time"
)

/*
 * Default state (USRENAME : PASSWORD):
 * 		test : test123
 * 		user : passwd
 * 		follower : follower
 */

func TestLogin(t *testing.T) {
	token := _Login("test", "test123")
	if len(token) == 0 {
		t.Errorf("Expect nonempty token, got empty")
	}
}

func TestSignup(t *testing.T) {
	tokenSignup := _Signup("newuser", "newpass")
	if len(tokenSignup) == 0 {
		t.Errorf("Expect nonempty signup token, got empty")
	}

	// wait for raft cluster to commit the update
	time.Sleep(300 * time.Millisecond)
	tokenLogin := _Login("newuser", "newpass")
	if len(tokenLogin) == 0 {
		t.Errorf("Expect nonempty login token, got empty")
	}
}
