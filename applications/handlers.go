package applications

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hipo/gotcha/mongo"
	"github.com/hipo/gotcha/users"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"
)

func ApplicationListTemplateHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("applications.tmpl")
	t = template.Must(t.ParseGlob("templates/*.tmpl"))
	t.Execute(w, nil)
}

func UrlListTemplateHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("urls.tmpl")
	t = template.Must(t.ParseGlob("templates/*.tmpl"))
	application := Application{}
	vars := mux.Vars(r)

	applicationId := vars["applicationId"]
	applicationp := &application
	err := mongo.Find(application, bson.M{"_id": bson.ObjectIdHex(applicationId)}).One(applicationp)
	if err != nil {
		w.WriteHeader(404)
		return

	}
	t.Execute(w, applicationp)
}

func IsAuthenticated(authToken string, applicationId string) bool {
	user := users.User{}
	application := Application{}
	err := mongo.Find(user, bson.M{"token": authToken}).One(&user)
	if err != nil {
		panic(err)
	}
	err = mongo.Find(application, bson.M{"_id": bson.ObjectIdHex(applicationId)}).One(&application)
	if err != nil {
		panic(err)
	}
	if application.OwnerId != user.Id {
		return false
	}
	return true
}

func ApplicationListHandler(w http.ResponseWriter, r *http.Request) {

	var applications []Application

	user := users.User{}
	application := Application{}

	token := r.FormValue("token")
	err := mongo.Find(user, bson.M{"token": token}).One(&user)

	if err != nil {
		w.WriteHeader(403)
		return
	}
	err = mongo.Find(application, bson.M{"ownerid": user.Id}).All(&applications)

	if err != nil {
		w.WriteHeader(403)
		return
	}
	if err != nil {
		w.WriteHeader(400)
		return
	}
	json.NewEncoder(w).Encode(applications)
}

func ApplicationAddHandler(w http.ResponseWriter, r *http.Request) {
	application := Application{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	applicationp := &application
	err = json.Unmarshal(body, applicationp)
	if err != nil {
		panic(err)
	}
	user := users.User{}
	token := r.FormValue("token")
	err = mongo.Find(user, bson.M{"token": token}).One(&user)

	applicationp.Id = bson.NewObjectId()
	application.OwnerId = user.Id
	applicationp.CreateApplication()
	w.WriteHeader(200)
	return
}

func ApplicationDeleteHandler(w http.ResponseWriter, r *http.Request) {
	application := Application{}
	vars := mux.Vars(r)
	applicationId := vars["applicationId"]
	token := r.FormValue("token")

	isAuthenticated := IsAuthenticated(token, applicationId)
	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}

	err := mongo.Delete(application, bson.M{"_id": bson.ObjectIdHex(applicationId)})
	if err != nil {
		w.WriteHeader(400)
		fmt.Println(err)
		return
	}
	w.WriteHeader(204)
}

func UrlListHandler(w http.ResponseWriter, r *http.Request) {
	url := Url{}
	var urls []Url
	vars := mux.Vars(r)
	applicationId := vars["applicationId"]
	token := r.FormValue("token")

	isAuthenticated := IsAuthenticated(token, applicationId)
	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	err := mongo.Find(url, bson.M{"application_id": bson.ObjectIdHex(applicationId)}).All(&urls)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	serializedUrls := make([]map[string]interface{}, len(urls))

	for i, element := range urls {
		serializedUrls[i] = element.Serialize()
	}
	json.NewEncoder(w).Encode(serializedUrls)
}

func UrlAddHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	applicationId := vars["applicationId"]
	token := r.FormValue("token")
	isAuthenticated := IsAuthenticated(token, applicationId)

	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	urlp := &Url{}
	err = json.Unmarshal(body, urlp)
	if err != nil {
	}

	urlp.ApplicationId = bson.ObjectIdHex(applicationId)
	err = urlp.CreateUrl()

	if err != nil {
		w.WriteHeader(400)
		return
	}
}

func UrlDeleteHandler(w http.ResponseWriter, r *http.Request) {

	url := Url{}
	vars := mux.Vars(r)
	applicationId := vars["applicationId"]
	urlId := vars["urlId"]

	token := r.FormValue("token")
	isAuthenticated := IsAuthenticated(token, applicationId)

	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	err := mongo.Delete(url, bson.M{"_id": bson.ObjectIdHex(urlId),
		"application_id": bson.ObjectIdHex(applicationId)})

	if err != nil {
		w.WriteHeader(400)
		fmt.Println(err)
		return
	}
	w.WriteHeader(204)

}

func FetchURL(url string, UrlId bson.ObjectId) {
	time_start := time.Now()
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	urlRecord := UrlRecord{}

	urlRecordp := &urlRecord
	urlRecordp.Id = bson.NewObjectId()
	urlRecordp.UrlId = UrlId
	urlRecordp.Time = time.Since(time_start).Seconds()
	urlRecordp.CreateUrlRecord()
}

func FetchApplicationURLs(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	applicationId := vars["applicationId"]
	token := r.FormValue("token")
	isAuthenticated := IsAuthenticated(token, applicationId)

	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	url := Url{}
	var urls []Url
	err := mongo.Find(url, bson.M{"application_id": bson.ObjectIdHex(applicationId)}).All(&urls)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	for i := 0; i < len(urls); i++ {
		go FetchURL(urls[i].Url, urls[i].Id)
	}
}
