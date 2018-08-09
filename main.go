package main

import (
	. "github.com/phips4/communityTools/app"
	"github.com/phips4/communityTools/app/db"
	"github.com/phips4/communityTools/handler"
	"github.com/phips4/communityTools/server"
	"log"
)

func main() {
	//logger
	defer loadLogger()()

	//config
	conf := loadConfig()
	_ = conf

	// mgo
	mgoSession := db.ConnectMongo(conf.MongoDB.Host, conf.MongoDB.Port)
	defer mgoSession.Close()

	log.Println("starting application")

	//webserver
	webServer := server.New(mgoSession)
	webServer.Init()
	//register all handlers
	handler.AddAllPollHandler(webServer)
	webServer.Run()

}

func loadLogger() func() {
	logger := NewLogWriter()
	log.SetFlags(0)
	log.SetOutput(logger)
	log.Println("logger initialized.")

	return func() {
		logger.Close()
	}
}

func loadConfig() *ConfigStruct {
	config := Config{}
	log.Println("config read.")
	return config.LoadConfig()
}
