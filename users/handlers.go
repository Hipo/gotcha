package users

import (
	"encoding/json"
	"github.com/hipo/gotcha/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"text/template"
)

func LoginTemplateHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("login.tmpl")
	t = template.Must(t.ParseGlob("templates/*.tmpl"))
	t.Execute(w, nil)
}

func SignupTemplateHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("signup.tmpl")
	t = template.Must(t.ParseGlob("templates/*.tmpl"))
	t.Execute(w, nil)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	user := User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	userp := &user
	err = json.Unmarshal(body, userp)
	err = userp.CreateUser()
	if err != nil {
		w.WriteHeader(400)
	}
	err = mongo.Find(user, bson.M{"email": userp.Email}).One(userp)
	json.NewEncoder(w).Encode(userp.Serialize())

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	user := User{}
	requestuser := User{}
	err = json.Unmarshal(body, &requestuser)
	err = mongo.Find(user, bson.M{"email": requestuser.Email}).One(&user)

	if err != nil {
		w.WriteHeader(403)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestuser.Password))

	if err != nil {
		w.WriteHeader(403)
		return
	}
	userp := &user
	json.NewEncoder(w).Encode(userp.Serialize())

}
