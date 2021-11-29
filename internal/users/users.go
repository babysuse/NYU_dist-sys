package users

import (
	"database/sql"
	"log"

	database "github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/db/migrations/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"ID"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) Create() {
	statement, err := database.DB.Prepare("INSERT INTO Users(Username, Password) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := HashPassword(user.Password)
	_, err = statement.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Println(err)
		return
	}
}

func (user *User) Authenticate() bool {
	statement, err := database.DB.Prepare("SELECT Password FROM Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}

	row := statement.QueryRow(user.Username)
	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return CheckPasswordHash(user.Password, hashedPassword)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.DB.Prepare("SELECT ID FROM Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var id int
	err = row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return 0, err
	}
	return id, nil
}

func (user *User) Follow(followeeID string) {
	stmt, err := database.DB.Prepare("INSERT INTO Following(UserID, FolloweeID) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(user.ID, followeeID)
	if err != nil {
		log.Println(err)
		return
	}
}

func (user *User) Unfollow(followeeID string) {
	stmt, err := database.DB.Prepare("DELETE FROM Following WHERE UserID = ? AND FolloweeID = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(user.ID, followeeID)
	if err != nil {
		log.Println(err)
		return
	}
}

func (user *User) GetFollowee() []User {
	rows, err := database.DB.Query(`
		SELECT U.ID, U.Username
		FROM Following F join Users U ON F.FolloweeID = U.ID
		WHERE F.UserID = ?`, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var followees []User
	for rows.Next() {
		var followee User
		err := rows.Scan(&followee.ID, &followee.Username)
		if err != nil {
			log.Fatal(err)
		}
		followees = append(followees, followee)
	}
	return followees
}
