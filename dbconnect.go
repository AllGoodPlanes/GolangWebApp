package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
	"log"
)

const (
	MongoDBHosts = "xxxxxxxxx"
	AuthDatabase = "xxxxxxxxxxxx"
	AuthUserName = "xxxxxxxx"
	AuthPassword = "xxxxxxxxxxx"
	TestDatabase = "xxxxxxxxx"
)

var db *sql.DB
var mongoSession *mgo.Session

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	mongoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatal(err)
	}

}
