package db

import (
	"fmt"
	"github.com/phips4/communityTools/app"
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

func Connect(mgoConf app.MgoConf) *mgo.Session {
	if session != nil {
		session.Close()
	}
	var ses *mgo.Session
	var err error
	auth := func() {
		if mgoConf.User == "" {
			ses, err = mgo.Dial(fmt.Sprintf("%s:%d", mgoConf.Host, mgoConf.Port))
		} else {
			ses, err = mgo.DialWithInfo(&mgo.DialInfo{
				Database: mgoConf.Database,
				Username: mgoConf.User,
				Password: mgoConf.Password,
				Addrs:    []string{fmt.Sprintf("%s:%d", mgoConf.Host, mgoConf.Port)},
			})
		}
	}
	if err != nil {
		log.Fatal(err)
		return nil
	}
	auth()
	wait := 5

	for err != nil {
		log.Printf("can not connect to mongodb (%s:%d: error: %s) Waiting %d secounds.", mgoConf.Host, mgoConf.Port, err, wait)
		time.Sleep(time.Second * time.Duration(wait))
		auth()
		wait += 5
	}

	log.Println("successfully connected to mongodb.")

	session = ses
	return ses
}

/* ======================================
 * | PollSession                        |
 * ====================================== */
type PollSession struct {
	c        *mgo.Collection
	closePtr func()
}

func GetPollSession() *PollSession {
	s := session.Clone()
	return &PollSession{s.DB(dbName).C(pollCollection), s.Close}
}

func (ps *PollSession) Close() {
	ps.closePtr()
}

func (ps *PollSession) PollExists(id string) (bool, error) {
	check := struct {
		ID string `bson:"_id"`
	}{}

	err := ps.c.FindId(id).One(&check)
	if err != nil {
		// it is safer to say the ID already exists, it would not create a new one
		return true, err
	}

	return check.ID == id, nil
}

func (ps *PollSession) InsertPoll(doc interface{}) error {
	return ps.c.Insert(doc)
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
