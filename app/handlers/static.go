package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/app/servers"
	"io/ioutil"
)

//TODO: add frontend git submodule instead or consider to put the frontend as its own module instead of the defaultServer
func AddAllStaticRoutes(server *servers.DefaultServer) {
	server.Router.StaticFile("/", "./public/static/index.html")

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
