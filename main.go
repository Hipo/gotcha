package main

import (
	"github.com/hipo/gotcha/mongo"
	"log"
	"net/http"
	"fmt"
	"database/sql"
)

func main() {
	router := NewRouter()
	mongo.Connect(Config.DbHostString(), Config.DbName)
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        Config.Username, Config.Password, Config.DbName)
    	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Fatal(http.ListenAndServe(":8080", router))
}
