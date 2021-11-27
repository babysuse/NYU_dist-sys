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
	User users.User
}

func (p Post) Save() int64 {
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
