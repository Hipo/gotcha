package applications

import (
	"fmt"
	"github.com/hipo/gotcha/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Application struct {
	Id          int	 	  `json:"id" db:"_id"`
	OwnerId     int 	  `json:"owner_id" db:"owner_id"`
	CallbackUrl string        `json:"callback_url" db:"callback_url"`
	WaitTime    int           `json:"wait_time" db:"callback_url"`
	Name        string        `json:"name" db:"name"`
}

func (a Application) Collection() string { return "applications" }
func (a Application) CreateApplication() error {
	err := mongo.Insert(a)
	if err != nil {
		return err
	}
	return nil
}

func (a Application) UrlCount() string {
	url := Url{}
	count, err := mongo.Find(url, bson.M{"application_id": a.Id}).Count()
	if err != nil {
		fmt.Println(err)

	}
	return fmt.Sprintf("%v", count)
}

func (a Application) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"Name":        a.Name,
		"Id":          a.Id,
		"Count":       a.UrlCount(),
		"CallbackUrl": a.CallbackUrl,
	}
}

type Url struct {
	Id            int		`json:"id" db:"_id"`
	Url           string            `json:"url" db:"url"`
	Title         string            `json:"title" db:"title"`
	WaitTime      int               `json:"wait_time" db:"wait_time"`
	TryCount      int               `json:"try_count" db:"try_count"`
	ApplicationId int 		`json:"application_id" db:"application_id"`
	Headers       map[string]string `json:"headers" db:"headers"`
}

func (u Url) Collection() string { return "urls" }

func (u Url) CreateUrl() error {
	u.Id = bson.NewObjectId()
	mongo.Insert(u)
	return nil
}

func (u Url) UpdateUrl(filter_query interface{}, update_data interface{}) error {
	err := mongo.Update(u, filter_query, update_data)
	if err != nil {
		return err
	}
	return err
}

func (u *Url) Deserialize() (map[string] interface{}) {
	bundle := make(map[string]interface{})
	bundle["url"] = u.Url
	bundle["title"] = u.Title
	bundle["wait_time"] = u.WaitTime
	bundle["try_count"] = u.TryCount
	bundle["application_id"] = u.ApplicationId
	bundle["headers"] = u.Headers
	return bundle
}


func (u *Url) Serialize() (map[string]interface{}, error) {
	record := UrlRecord{}
	records := []UrlRecord{}
		err := mongo.Find(record, bson.M{"url_id": bson.ObjectId(u.Id)}).Sort("-date_created").Limit(2).All(&records)
	if err != nil {
		return nil, err
	}
	bundle := make(map[string]interface{})
	bundle["Id"] = u.Id
	bundle["Url"] = u.Url
	bundle["Title"] = u.Title
	bundle["WaitTime"] = u.WaitTime
	bundle["TryCount"] = u.TryCount
	bundle["ApplicationId"] = u.ApplicationId
	bundle["Headers"] = u.Headers
	if len(records) >= 1 {
		record1 := records[0]
		bundle["Time"] = record1.DateCreated
		bundle["Last"] = record1.Time
		bundle["Status"] = record1.StatusCode
	}

	if len(records) >= 2 {

		record1, record2 := records[0], records[1]
		bundle["Previous"] = record2.Time
		bundle["Faster"] = record1.Time < record2.Time
	}

	return bundle, nil
}

type UrlRecord struct {
	Id          int 	  `json:"id" db:"_id"`
	UrlId       int 	  `json:"url_id" db:"url_id"`
	Time        float64       `json:"time" db:"time"`
	StatusCode  string        `json:"status_code" db:"status_code"`
	DateCreated time.Time     `json:"date_created" db:"date_created"`
}

func (r *UrlRecord) Serialize() map[string]interface{} {
	bundle := make(map[string]interface{})
	bundle["Time"] = r.Time
	bundle["StatusCode"] = r.StatusCode
	bundle["DateCreated"] = r.DateCreated

	return bundle
}

func (u UrlRecord) Collection() string { return "urlrecords" }
func (u UrlRecord) CreateUrlRecord() error {

	u.DateCreated = time.Now()
	err := mongo.Insert(u)
	if err != nil {
		panic(err)
	}
	return nil
}
