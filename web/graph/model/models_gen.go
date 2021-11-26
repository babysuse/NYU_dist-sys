// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreatePost struct {
	Text     string `json:"text"`
	AuthorID string `json:"authorId"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Post struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Author *User  `json:"author"`
}

type RefreshToken struct {
	Token string `json:"token"`
}

type Signup struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
