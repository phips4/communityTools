package main

import (
	"fmt"
	. "github.com/phips4/communityTools/app"
	"github.com/phips4/communityTools/app/db"
	"github.com/phips4/communityTools/handler"
	"github.com/phips4/communityTools/server"
	"log"
)

func main() {

	println("#==============================================================================#")
	println(`|                                            _ _      _______          _       |
|                                           (_) |    |__   __|        | |      |
|   ___ ___  _ __ ___  _ __ ___  _   _ _ __  _| |_ _   _| | ___   ___ | |___   |
|  / __/ _ \| '_ ' _ \| '_ ' _ \| | | | '_ \| | __| | | | |/ _ \ / _ \| / __|  |
| | (_| (_) | | | | | | | | | | | |_| | | | | | |_| |_| | | (_) | (_) | \__ \  |
|  \___\___/|_| |_| |_|_| |_| |_|\__,_|_| |_|_|\__|\__, |_|\___/ \___/|_|___/  |
|                                                   __/ |                      |
|                                                  |___/                       |`)
	println("#==============================================================================#")

	//logger
	defer loadLogger()()

	//config
	conf := loadConfig()

	// mgo
	mongo := db.Connect(conf.MongoDB.Host, conf.MongoDB.Port)
	defer mongo.Close()

	db.Login(conf.MongoDB.Database, nil)

	log.Println("starting application")

	//webserver
	webServer := server.New()
	//register all handlers
	handler.AddAllPollHandler(webServer)
	webServer.Listen(fmt.Sprintf("%s:%d", conf.WebServer.Host, conf.WebServer.Port))
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
