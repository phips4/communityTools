package server

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

type WebServer struct {
	Router     *gin.Engine
	MgoSession *mgo.Session
}

func New(mgoSession *mgo.Session) *WebServer {
	return &WebServer{
		nil,
		mgoSession,
	}
}

func (server *WebServer) Init() {
	server.Router = gin.Default()
}

func (server *WebServer) Run() {
	server.Router.Run(":4337")
}

func (server *WebServer) Stop() {
	//TODO: graceful shutdown
}
