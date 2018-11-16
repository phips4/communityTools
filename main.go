package main

import (
	"context"
	"github.com/gin-gonic/gin"
	. "github.com/phips4/communityTools/app"
	"github.com/phips4/communityTools/app/db"
	"github.com/phips4/communityTools/app/handlers"
	"github.com/phips4/communityTools/app/servers"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	println(
`#==============================================================================#
|                                            _ _      _______          _       |
|                                           (_) |    |__   __|        | |      |
|   ___ ___  _ __ ___  _ __ ___  _   _ _ __  _| |_ _   _| | ___   ___ | |___   |
|  / __/ _ \| '_ ' _ \| '_ ' _ \| | | | '_ \| | __| | | | |/ _ \ / _ \| / __|  |
| | (_| (_) | | | | | | | | | | | |_| | | | | | |_| |_| | | (_) | (_) | \__ \  |
|  \___\___/|_| |_| |_|_| |_| |_|\__,_|_| |_|_|\__|\__, |_|\___/ \___/|_|___/  |
|                                                   __/ |                      |
|                                                  |___/                       |
#==============================================================================#`)

	//logger
	closeLogger := loadLogger()

	if os.Getenv("GIN_MODE") == "" {
		log.Println("[INFO]: set env. 'export GIN_MODE=release' on production")
		gin.SetMode(gin.DebugMode)
	}

	//config
	conf := loadConfig()

	// mgo
	mongo := db.Connect(conf.MgoConf)
	defer mongo.Close()

	log.Println("starting application")

	// default webserver
	var dSrv *servers.DefaultServer
	dSrv = servers.NewDefaultServer()
	handlers.AddAllStaticRoutes(dSrv)
	handlers.AddAllPollHandler(dSrv)
	handlers.AddAllGeneralHandler(dSrv)

	// start default server
	go dSrv.Listen(conf.Webserver.Host, conf.Webserver.Port)

	waitForShutdown(func() {
		//close default webserver
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := dSrv.Stop(ctx); err != nil {
			log.Fatal("Server Shutdown: ", err)
		}

		closeLogger()
	})
}

func waitForShutdown(shutdownEvent func()) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown app...")
	shutdownEvent()
	log.Println("good bye!")
}

func loadLogger() func() {
	logger := NewLogWriter()
	log.SetFlags(0)
	log.SetOutput(logger)
	log.Println("logger initialized.")

	return func() {
		logger.Close()
		err := logger.Compress()
		if err != nil {
			panic(err)
		}
	}
}

func loadConfig() *ConfigStruct {
	config := Config{}
	log.Println("config read.")
	return config.LoadConfig()
}