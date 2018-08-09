package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

func ConnectMongo(host string, port int) *mgo.Session {
	session, err := mgo.Dial(fmt.Sprintf("%s:%d", host, port))
	wait := 5

	for err != nil {
		log.Printf("can not connect to mongodb (%s:%d: %s) error: Waiting %d secounds.", host, port, err, wait)
		time.Sleep(time.Second * time.Duration(wait))
		session, err = mgo.Dial(fmt.Sprintf("%s:%d", host, port))
		wait += 5
	}

	log.Println("successfully connected to mongodb.")

	return session
}
