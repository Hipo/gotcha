package main

import (
	"github.com/hipo/gotcha/mongo"
	"log"
	"net/http"
)

func main() {
	router := NewRouter()
	mongo.Connect(Config.DbHostString(), Config.DbName)
	log.Fatal(http.ListenAndServe(":8080", router))
}
