package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/app/servers"
	"github.com/phips4/communityTools/app/utils"
	"log"
	"net/http"
	"strings"
)

//TODO: move to config file
var MODULES = []string {"image.", "txt.", "polls.", "file."}

// subdomain specific routing
func indexGET(ctx *gin.Context) (int, error) {
	subdomain := ""
	if utils.HasPrefix(ctx.Request.Host, MODULES) {
		subdomain = strings.Split(ctx.Request.Host, ".")[0]
	}

	ctx.String(http.StatusOK, "module: %s", subdomain)
	ctx.Abort()

	return http.StatusOK, nil
}

func AddAllGeneralHandler(server *servers.DefaultServer) {
	if gin.IsDebugging() {
		log.Print(" ")
		log.Print("GENERAL HANDLERS")
	}
	server.Router.GET("/", errorHandler(indexGET))
}