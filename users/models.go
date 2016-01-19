package users

import (
	"crypto/rand"
	"fmt"
	"github.com/hipo/gotcha/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func randToken() string {
	b := make([]byte, 12)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type User struct {
	Id       bson.ObjectId `json: "id" bson:"_id"`
	Email    string        `json: "email" bson: "email"`
	Password string        `json: "password" bson: "password"`
	Token    string        `json: "token" bson: "token"`
}

func (u User) Collection() string { return "users" }
func (u User) CreateUser() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
	u.Token = randToken()
	u.Id = bson.NewObjectId()
	err = mongo.Insert(u)
	fmt.Println(err)
	return err
}
func (u *User) Serialize() map[string]string {
	return map[string]string{
		"email": u.Email,
		"id":    u.Id.String(),
		"token": u.Token,
	}
}
