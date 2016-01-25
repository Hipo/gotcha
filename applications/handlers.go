package applications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hipo/gotcha/mongo"
	"github.com/hipo/gotcha/users"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"math"
	"net/http"
	"text/template"
	"time"
)

type Callback struct {
	CallbackUrl string `json:"callbackurl"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.tmpl")
	t = template.Must(t.ParseGlob("templates/*.tmpl"))
	t.Execute(w, nil)
}

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

func UrlDetailTemplateHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("url_detail.tmpl")
	t = template.Must(t.ParseGlob("templates/*.tmpl"))
	url := Url{}
	vars := mux.Vars(r)

	urlId := vars["urlId"]
	urlp := &url
	err := mongo.Find(url, bson.M{"_id": bson.ObjectIdHex(urlId)}).One(urlp)
	if err != nil {
		w.WriteHeader(404)
		return

	}
	t.Execute(w, urlp)
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
	serializedApplications := make([]map[string]string, len(applications))

	for i, element := range applications {
		serializedApplications[i] = element.Serialize()
	}
	json.NewEncoder(w).Encode(serializedApplications)
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

func AddCallbackHandler(w http.ResponseWriter, r *http.Request) {
	application := Application{}
	vars := mux.Vars(r)
	applicationId := vars["applicationId"]
	token := r.FormValue("token")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(400)
	}

	isAuthenticated := IsAuthenticated(token, applicationId)
	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	callback := Callback{}
	json.Unmarshal(body, &callback)
	mongo.Update(application,
		bson.M{"_id": bson.ObjectIdHex(applicationId)},
		bson.M{"$set": bson.M{"callbackurl": callback.CallbackUrl}})

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

	err = json.NewEncoder(w).Encode(serializedUrls)
}

func UrlAddHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
func UrlDetailHandler(w http.ResponseWriter, r *http.Request) {

	url_record := UrlRecord{}
	var records []UrlRecord
	vars := mux.Vars(r)
	urlId := vars["urlId"]
	applicationId := vars["applicationId"]
	token := r.FormValue("token")

	isAuthenticated := IsAuthenticated(token, applicationId)
	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	err := mongo.Find(url_record, bson.M{"url_id": bson.ObjectIdHex(urlId)}).All(&records)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	serializedUrlRecords := make([]map[string]interface{}, len(records))

	for i, element := range records {
		serializedUrlRecords[i] = element.Serialize()
	}
	json.NewEncoder(w).Encode(serializedUrlRecords)
}

func AvarageAccordingStandardDeviation(timelist []float64, mean float64, deviation float64) float64 {
	max := mean + deviation
	min := mean - deviation

	count := 0.0
	sum := 0.0

	for i := range timelist {
		if (timelist[i] < max) || (timelist[i] > min) {
			count = count + 1
			sum = sum + timelist[i]
		}
	}
	avarage := sum / count
	return avarage
}

func CalculateStandardDeviation(timelist []float64) (float64, float64) {
	total := 0.0
	variance := 0.0
	for i := 0; i < len(timelist); i++ {
		total = total + timelist[i]
	}

	mean := total / float64(len(timelist))
	for i := 0; i < len(timelist); i++ {
		variance += math.Pow((timelist[i] - mean), 2.0)
	}
	variance = variance / (float64(len(timelist)) - 1.0)
	deviation := math.Sqrt(variance)
	return deviation, mean
}

func FetchThread(url Url, timelist chan float64, statusList chan string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url.Url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range url.Headers {
		req.Header.Set(k, v)
	}
	time_start := time.Now()
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Couldn't connect to server")
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	timeSpent := time.Since(time_start).Seconds()
	timelist <- timeSpent
	fmt.Println(timeSpent)
	statusList <- response.Status
}

func FetchURL(channel chan bool, url Url, UrlId bson.ObjectId) {
	statusCode := ""
	tryCount := 3
	timelist := make(chan float64)
	statusList := make(chan string)

	times := make([]float64, tryCount)

	for i := 0; i < tryCount; i++ {
		time.Sleep(2 * time.Second)
		go FetchThread(url, timelist, statusList)
	}
	for i := 0; i < tryCount; i++ {
		time := <-timelist
		times[i] = time
	}
	statusCode = <-statusList
	deviation, mean := CalculateStandardDeviation(times)
	avarageTime := AvarageAccordingStandardDeviation(times, mean, deviation)
	urlRecord := UrlRecord{}
	urlRecordp := &urlRecord
	urlRecordp.Id = bson.NewObjectId()
	urlRecordp.UrlId = UrlId
	urlRecordp.StatusCode = statusCode
	urlRecordp.Time = avarageTime
	urlRecordp.CreateUrlRecord()
	channel <- true
	return
}

func PostCallback(channel chan bool, count int, url string, applicationId string) {
	_url := Url{}
	var urls []Url

	for i := 0; i < count; i++ {
		<-channel
	}
	mongo.Find(_url, bson.M{"application_id": bson.ObjectIdHex(applicationId)}).All(&urls)
	urlList := make([]map[string]interface{}, len(urls))

	for i := 0; i < len(urls); i++ {
		urlList = append(urlList, urls[i].Serialize())
	}
	urlJSON, _ := json.Marshal(urlList)
	urlJSONstring := string(urlJSON)
	var jsonUrls = []byte(urlJSONstring)
	http.Post(url, "application/json", bytes.NewBuffer(jsonUrls))
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
	application := Application{}
	var urls []Url
	err := mongo.Find(url, bson.M{"application_id": bson.ObjectIdHex(applicationId)}).All(&urls)
	err = mongo.Find(application, bson.M{"_id": bson.ObjectIdHex(applicationId)}).One(&application)

	if err != nil {
		w.WriteHeader(400)
		return
	}
	channel := make(chan bool)

	for i := 0; i < len(urls); i++ {
		go FetchURL(channel, urls[i], urls[i].Id)
	}

	go PostCallback(channel, len(urls), application.CallbackUrl, applicationId)

}
