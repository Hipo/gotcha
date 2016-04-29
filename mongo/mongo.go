package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Connection struct {
	Session *mgo.Session
	Db      *mgo.Database
	Url     string
	Host    string
	Port    string
}

type Model interface {
	Collection() string
}

func Cursor(m Model) *mgo.Collection {
	return Current.Db.C(m.Collection())
}

func Find(m Model, query interface{}) *mgo.Query {
	cursor := Cursor(m)
	return cursor.Find(query)
}

func FindId(m Model, id bson.ObjectId) *mgo.Query {
	cursor := Cursor(m)
	return cursor.FindId(id)
}

func Insert(m Model) (err error) {
	return Current.Db.C(m.Collection()).Insert(m)
}

func Update(m Model,
	filter_query interface{},
	update_query interface{}) (err error) {
	return Current.Db.C(m.Collection()).Update(filter_query, update_query)
}

func Delete(m Model,
	delete_query interface{}) (err error) {
	return Current.Db.C(m.Collection()).Remove(delete_query)
}

var Current = new(Connection)

func Connect(url, database string) {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	db := session.DB(database)
	Current.Session = session
	Current.Db = db
}
