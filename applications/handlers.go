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
	"sync"
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

func IsAuthenticated(authToken string, applicationId string) (bool, error) {
	user := users.User{}
	application := Application{}
	err := mongo.Find(user, bson.M{"token": authToken}).One(&user)
	if err != nil {
		return false, err
	}
	err = mongo.Find(application, bson.M{"_id": bson.ObjectIdHex(applicationId)}).One(&application)
	if err != nil {
		return false, err
	}
	if application.OwnerId != user.Id {
		return false, nil
	}
	return true, nil
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
	serializedApplications := make([]map[string]interface{}, len(applications))

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

	isAuthenticated, err := IsAuthenticated(token, applicationId)
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

	isAuthenticated, err := IsAuthenticated(token, applicationId)
	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}

	err = mongo.Delete(application, bson.M{"_id": bson.ObjectIdHex(applicationId)})
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

	isAuthenticated, err := IsAuthenticated(token, applicationId)
	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}

	err = mongo.Find(url, bson.M{"application_id": bson.ObjectIdHex(applicationId)}).All(&urls)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	serializedUrls := make([]map[string]interface{}, len(urls))
	for i, element := range urls {
		serializedUrls[i], err = element.Serialize()
	}
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}


	err = json.NewEncoder(w).Encode(serializedUrls)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

}

func UrlAddHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	applicationId := vars["applicationId"]
	token := r.FormValue("token")

	isAuthenticated, err := IsAuthenticated(token, applicationId)

	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	urlp := &Url{WaitTime:100,
		    TryCount:10}
	err = json.Unmarshal(body, urlp)
	if err != nil {
		fmt.Println(err)
	}

	urlp.ApplicationId = bson.ObjectIdHex(applicationId)


	err = urlp.CreateUrl()

	if err != nil {
		w.WriteHeader(400)
		return
	}
}

func UrlEditHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	urlId := vars["urlId"]
	var url Url
	err := mongo.Find(url, bson.M{"_id": bson.ObjectIdHex(urlId)}).One(&url)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	applicationId := vars["applicationId"]
	token := r.FormValue("token")
	isAuthenticated, err := IsAuthenticated(token, applicationId)
	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	urlp := &Url{WaitTime:100,
		     TryCount:10}

	err = json.Unmarshal(body, urlp)
	if err != nil {
		fmt.Println(err)
	}

	urlp.ApplicationId = bson.ObjectIdHex(applicationId)
	urlQuery := bson.M{"_id": url.Id}
	updateData := urlp.Deserialize()

	err = urlp.UpdateUrl(urlQuery, updateData)
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
	isAuthenticated, err := IsAuthenticated(token, applicationId)

	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	err = mongo.Delete(url, bson.M{"_id": bson.ObjectIdHex(urlId),
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
	url := Url{}
	vars := mux.Vars(r)
	urlId := vars["urlId"]
	applicationId := vars["applicationId"]
	token := r.FormValue("token")

	isAuthenticated, err := IsAuthenticated(token, applicationId)
	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	err = mongo.Find(url_record, bson.M{"url_id": bson.ObjectIdHex(urlId)}).All(&records)
	err = mongo.Find(url, bson.M{"_id": bson.ObjectIdHex(urlId)}).One(&url)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	serializedUrlRecords := make([]map[string]interface{}, len(records))

	for i, element := range records {
		serializedUrlRecords[i] = element.Serialize()
	}
	serializedUrl, err := url.Serialize()
	bundle := map[string]interface{}{"records": serializedUrlRecords,
		                         "url": serializedUrl}
	json.NewEncoder(w).Encode(bundle)
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
	if (count > 0) {
		avarage := sum / count
		return avarage
	}
	return 0.0
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


func FetchThread(url Url, timelist chan float64, statusList chan string, tryWait *sync.WaitGroup) {
	client := &http.Client{Timeout:time.Duration(30 * time.Second)}
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
		timelist <- 100
		statusList <- "503"
		return
	}
	defer response.Body.Close()
	timeSpent := time.Since(time_start).Seconds()
	tryWait.Done()
	timelist <- timeSpent
	statusList <- response.Status

}

func FetchURL(url Url, UrlId bson.ObjectId, wg *sync.WaitGroup) {

	statusCode := ""
	tryCount := url.TryCount
	if tryCount == 0 {
		tryCount = 1
	}

	timelist := make(chan float64)
	statusList := make(chan string)
	var tryWait sync.WaitGroup

	times := make([]float64, tryCount)

	for i := 0; i < tryCount; i++ {
		time.Sleep(time.Duration(url.WaitTime) * time.Millisecond)
		tryWait.Add(1)
		go FetchThread(url, timelist, statusList, &tryWait)
		tryWait.Wait()
	}
	for i := 0; i < tryCount; i++ {
		time := <-timelist
		times[i] = time
	}
	statusCode = <-statusList

	avarageTime := times[0]
	if len(times) > 1 {
		deviation, mean := CalculateStandardDeviation(times)
		avarageTime = AvarageAccordingStandardDeviation(times, mean, deviation)
	}
	urlRecord := UrlRecord{Time: 100, StatusCode: ""}
	urlRecordp := &urlRecord
	urlRecordp.Id = bson.NewObjectId()
	urlRecordp.UrlId = UrlId
	urlRecordp.StatusCode = statusCode
	urlRecordp.Time = avarageTime
	urlRecordp.CreateUrlRecord()
	wg.Done()
	return
}

func PostCallback(count int, url string, applicationId string) {

	applicationUrl := "http://gotcha.hipo.biz/applications/" + applicationId + "/urls"
	message := fmt.Sprintf("Benchmark completed like a boss! <%s|Check details.>", applicationUrl)
	callbackData := map[string]string{"username": "gotcha",
					  "text": message}


	callbackDataJSON, _ := json.Marshal(callbackData)
	callbackDataString := string(callbackDataJSON)
	callbackDataB := []byte(callbackDataString)

	http.Post(url, "application/json", bytes.NewBuffer(callbackDataB))
}

func AsyncUrlCall(urls []Url, application Application, applicationId string, postCallback bool){
	var wg sync.WaitGroup
	for i := 0; i < len(urls); i++ {
	        wg.Add(1)
		go FetchURL(urls[i], urls[i].Id, &wg)
		wg.Wait()
	}
	if postCallback == true {
		go PostCallback(len(urls), application.CallbackUrl, applicationId)
	}
}

func FetchURLHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	token := r.FormValue("token")
	applicationId := vars["applicationId"]
	urlId := vars["urlId"]
	isAuthenticated, err := IsAuthenticated(token, applicationId)

	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	url := Url{}
	application := Application{}
	var urls []Url
	err = mongo.Find(application, bson.M{"_id": bson.ObjectIdHex(applicationId)}).One(&application)

	err = mongo.Find(url, bson.M{"application_id": bson.ObjectIdHex(applicationId),
				     "_id": bson.ObjectIdHex(urlId)}).All(&urls)

	if err != nil {
		w.WriteHeader(400)
		return
	}
	go AsyncUrlCall(urls, application, applicationId, false)

}


func FetchApplicationURLs(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	applicationId := vars["applicationId"]
	token := r.FormValue("token")
	isAuthenticated, err := IsAuthenticated(token, applicationId)

	if isAuthenticated != true {
		w.WriteHeader(403)
		return
	}
	url := Url{}
	application := Application{}
	var urls []Url
	err = mongo.Find(url, bson.M{"application_id": bson.ObjectIdHex(applicationId)}).All(&urls)
	err = mongo.Find(application, bson.M{"_id": bson.ObjectIdHex(applicationId)}).One(&application)

	if err != nil {
		w.WriteHeader(400)
		return
	}
	go AsyncUrlCall(urls, application, applicationId, true)


}
