package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/app"
	"github.com/phips4/communityTools/app/servers"
	"github.com/phips4/communityTools/app/utils"
	"net/http"
	"reflect"
	"strings"
)

//TODO: Frontend stuff

//TODO: move to config file
var MODULES []string
var conf *app.ConfigStruct

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

func modulesGET(ctx *gin.Context) (int, error) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "ok", "msg": "ok", "modules": conf.Modules})
	return http.StatusOK, nil
}

func AddAllGeneralHandler(server *servers.DefaultServer, config *app.ConfigStruct) {
	conf = config
	debugHandlerRegistration("GENERAL")

	initModuleNames()
	server.Router.GET("/", errorHandler(indexGET))
	server.Router.GET("/modules", errorHandler(modulesGET))
}

func initModuleNames() {
	ref := reflect.ValueOf(conf.Modules)
	modules := ref.NumField()
	MODULES = make([]string, modules)

	for i := 0; i < modules; i++ {
		MODULES[i] = ref.Field(i).FieldByName("ShortName").String()
	}
}