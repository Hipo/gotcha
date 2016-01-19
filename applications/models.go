package applications

import (
	"github.com/hipo/gotcha/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Application struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	OwnerId bson.ObjectId `json: "owner_id" bson:"owner_id"`
	Name    string        `json:"name" bson:"name"`
}

func (a Application) Collection() string { return "applications" }
func (a Application) CreateApplication() error {
	err := mongo.Insert(a)
	if err != nil {
		panic(err)
	}
	return nil
}

func (a *Application) Serialize() map[string]string {
	return map[string]string{
		"name": a.Name,
		"id":   a.Id.String(),
	}
}

type Url struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	Url           string        `json:"url" bson:"url"`
	Title         string        `json:"title" bson:"title"`
	ApplicationId bson.ObjectId `json:"application_id" bson:"application_id"`
}

func (u Url) Collection() string { return "urls" }
func (u Url) CreateUrl() error {
	u.Id = bson.NewObjectId()
	mongo.Insert(u)
	return nil
}

func (u *Url) Serialize() map[string]interface{} {
	record := UrlRecord{}
	records := []UrlRecord{}
	err := mongo.Find(record, bson.M{"url_id": bson.ObjectId(u.Id)}).Sort("-date_created").Limit(2).All(&records)
	if err != nil {
		panic(err)
	}
	if len(records) >= 2 {
		record1, record2 := records[0], records[1]
		return map[string]interface{}{
			"Id":       u.Id,
			"Url":      u.Url,
			"Title":    u.Title,
			"Last":     record1.Time,
			"Previous": record2.Time,
			"Time":     record1.DateCreated,
			"Faster":   record1.Time < record2.Time,
		}

	} else if len(records) == 1 {
		record1 := records[0]
		return map[string]interface{}{
			"Id":    u.Id,
			"Url":   u.Url,
			"Title": u.Title,
			"Last":  record1.Time,
		}
	} else {
		return map[string]interface{}{
			"Id":    u.Id,
			"Url":   u.Url,
			"Title": u.Title,
		}

	}

}

type UrlRecord struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	UrlId       bson.ObjectId `json:"url_id" bson:"url_id"`
	Time        float64       `json:"time" bson:"time"`
	DateCreated time.Time     `json:"date_created" bson:"date_created"`
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
