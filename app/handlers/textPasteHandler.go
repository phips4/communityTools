package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/app/servers"
)

func createPastePOST(ctx *gin.Context) (int, error) {
	return 0, nil
}

// register all paste endpoints
func AddAllPasteHandler(server *servers.DefaultServer) {
	debugHandlerRegistration("PASTE")

	pollGroup := server.Router.Group("/api/v1/paste")
	_ = pollGroup
	//pollGroup.POST("/create", errorHandler(createPastePOST))
	//pollGroup.GET("/get/:id", errorHandler(getPasteGET))
	//pollGroup.DELETE("/delete/:id", errorHandler(deletePasteDELETE))
}
