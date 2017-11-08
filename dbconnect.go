package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
	"log"
)

const (
	MongoDBHosts = "ds017678.mlab.com:17678"
	AuthDatabase = "heroku_549kcv64"
	AuthUserName = "heroku_549kcv64"
	AuthPassword = "s4viie11a5su1f0uj20r38i3qq"
	TestDatabase = "heroku_549kcv64"
)

var db *sql.DB
var mongoSession *mgo.Session

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://gzyvmlrlbuhxni:CFMYUehQv56yBGEWHdxy3psZ3T@ec2-54-195-241-96.eu-west-1.compute.amazonaws.com:5432/dcb0vrh3du4a7e")
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
