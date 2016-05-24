package users

import (
	"crypto/rand"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func randToken() string {
	b := make([]byte, 12)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type User struct {
	Id       int 	       `json:"id" db:"_id"`
	Email    string        `json:"email" db:"email"`
	Password string        `json:"password" db:"password"`
	Token    string        `json:"token" db:"token"`
}

func (u User) Collection() string { return "users" }
func (u User) CreateUser() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
        err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)

	fmt.Println(err)
	return err
}
func (u *User) Serialize() map[string]string {
	return map[string]string{
		"email": u.Email,
		"id":    u.Id,
		"token": u.Token,
	}
}
