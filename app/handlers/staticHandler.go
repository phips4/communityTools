package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/app/servers"
	"io/ioutil"
	"log"
)

//TODO: add frontend
func AddAllStaticRoutes(server *servers.DefaultServer) {
	if gin.IsDebugging() {
		log.Print(" ")
		log.Print("STATIC HANDLERS")
	}

	server.Router.NoRoute(func(context *gin.Context) {
		//only for testing, I know this is bad.
		bytes, err := ioutil.ReadFile("public/static/404.html")
		if err != nil {
			context.Error(err)
			return
		}
		context.Writer.Write(bytes)
		context.Done()
	})
}
