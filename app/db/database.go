package db

import (
	"fmt"
	"github.com/phips4/communityTools/app/polls"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

type SessionType int

const (
	pollCollection = "polls"
)

var (
	session *mgo.Session
	dbName  string
)

func Connect(host string, port int) *mgo.Session {

	if session != nil {
		session.Close()
	}

	ses, err := mgo.Dial(fmt.Sprintf("%s:%d", host, port))
	wait := 5

	for err != nil {
		log.Printf("can not connect to mongodb (%s:%d: %s) error: Waiting %d secounds.", host, port, err, wait)
		time.Sleep(time.Second * time.Duration(wait))
		ses, err = mgo.Dial(fmt.Sprintf("%s:%d", host, port))
		wait += 5
	}

	log.Println("successfully connected to mongodb.")

	session = ses
	return ses
}

func Login(database string, credential *mgo.Credential) error {
	dbName = database

	if credential == nil {
		return nil
	}

	return session.Login(credential)
}

func GetPollSession() *PollSession {
	s := session.Clone()
	return &PollSession{s.DB(dbName).C(pollCollection), s.Close}
}

/*
 * PollSession methods
 */
type PollSession struct {
	c            *mgo.Collection
	closePointer func()
}

func (ps *PollSession) Close() {
	ps.closePointer()
}

func (ps *PollSession) PollExists(id string) (bool, error) {
	type checkId struct {
		ID string `bson:"_id"`
	}

	var check *checkId
	err := ps.c.FindId(id).One(&check)

	if err != nil {
		// it is safer to say the ID already exists, because it would no create a new one
		// and we also check that case afterwards again
		return true, err
	}

	return check.ID == id, nil
}

func (ps *PollSession) InsertPoll(doc interface{}) error {
	return func() error {
		return ps.c.Insert(doc)
	}()
}

func (ps *PollSession) GetPoll(id string) (*polls.Poll, error) {
	var poll *polls.Poll
	err := ps.c.FindId(id).One(&poll)

	return poll, err
}

func (ps *PollSession) UpdatePoll(id string, poll *polls.Poll) error {
	return ps.c.UpdateId(id, poll)
}

func (ps *PollSession) DeletePoll(id string) error {
	return ps.c.RemoveId(id)
}