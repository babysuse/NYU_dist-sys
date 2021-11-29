package posts

import (
	"fmt"
	"log"

	database "github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/db/migrations/mysql"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/users"
)

type Post struct {
	ID   string
	Text string
	User *users.User
}

func (p Post) Save() int64 {
	println(p.Text, p.User.ID)
	stmt, err := database.DB.Prepare("INSERT INTO Posts(Text, AuthorID) Values (?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Exec(p.Text, p.User.ID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Post inserted!")
	return id

}

func GetAll(userID string) []Post {
	rows, err := database.DB.Query(`
		SELECT P.ID, P.text, U.ID, U.Username
		FROM Posts P join Users U on P.AuthorID = U.ID
		WHERE U.ID = ?
	`, userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var posts []Post
	var username string
	var id string
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Text, &id, &username)
		if err != nil {
			log.Fatal(err)
		}
		post.User = &users.User{
			ID:       id,
			Username: username,
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return posts
}
