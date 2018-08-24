package handler

import (
	"github.com/phips4/communityTools/server"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

//TODO: add frontend git submodule instead
func AddAllStaticRoutes(server *server.WebServer) {
	server.Router.StaticFile("/", "./public/static/index.html")

	server.Router.NoRoute(func(context *gin.Context) {
		//only for testing
		bytes, err := ioutil.ReadFile("public/static/404.html")
		if err != nil {
			context.Error(err)
			return
		}
		context.Writer.Write(bytes)
		context.Done()
	})
}