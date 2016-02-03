package applications

import (
	"fmt"
	"github.com/hipo/gotcha/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Application struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	OwnerId     bson.ObjectId `json: "owner_id" bson:"owner_id"`
	CallbackUrl string        `json: "callback_url" bson:"callbackurl"`
	Name        string        `json:"name" bson:"name"`
}

func (a Application) Collection() string { return "applications" }
func (a Application) CreateApplication() error {
	err := mongo.Insert(a)
	if err != nil {
		panic(err)
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

func (a Application) Serialize() map[string]string {
	return map[string]string{
		"Name":        a.Name,
		"Id":          a.Id.Hex(),
		"Count":       a.UrlCount(),
		"CallbackUrl": a.CallbackUrl,
	}
}

type Url struct {
	Id            bson.ObjectId     `json:"id" bson:"_id"`
	Url           string            `json:"url" bson:"url"`
	Title         string            `json:"title" bson:"title"`
	WaitTime      int               `json:"wait_time" bson:wait_time"`
	TryCount      int               `json:"try_count" bson:try_count"`
	ApplicationId bson.ObjectId     `json:"application_id" bson:"application_id"`
	Headers       map[string]string `json:"headers" bson:"headers"`
}

func (u Url) Collection() string { return "urls" }
func (u Url) CreateUrl() error {
	u.Id = bson.NewObjectId()
	mongo.Insert(u)
	return nil
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
	bundle["TryCount"] = u.WaitTime
	bundle["ApplicationId"] = u.ApplicationId
	if len(records) >= 1 {
		record1 := records[0]
		bundle["Last"] = record1.Time
		bundle["Time"] = record1.DateCreated
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
	Id          bson.ObjectId `json:"id" bson:"_id"`
	UrlId       bson.ObjectId `json:"url_id" bson:"url_id"`
	Time        float64       `json:"time" bson:"time"`
	StatusCode  string        `json:"status_code" bson:"status_code"`
	DateCreated time.Time     `json:"date_created" bson:"date_created"`
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
