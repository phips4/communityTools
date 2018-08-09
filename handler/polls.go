package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/server"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

func createPollPOST(ctx *gin.Context) {
	sess, ok := ctx.Get("mgo")
	defer sess.(*mgo.Session).Close()

	if !ok {
		AbortWithError(ctx, "db ctx doest not exist") //TODO: proper error handling
	}

	id := ctx.DefaultPostForm("id", "")
	/*
	title := ctx.DefaultPostForm("title", "")
	description := ctx.DefaultPostForm("description", "")
	cookieCheck := ctx.DefaultPostForm("cookieCheck", "false")
	multipleOptions := ctx.DefaultPostForm("multipleOptions", "false")
	options := ctx.PostFormArray("options")
	*/

	type idCheck struct {
		ID string `bson:"_id"`
	}

	var check idCheck
	err := sess.(*mgo.Session).DB("communityTools").C("polls").Find(bson.M{"_id": id}).One(&check)

	if err != nil {
		log.Println("error: " + err.Error())
		AbortWithError(ctx, err.Error()) //TODO: 500 Internal Server Error
	}

	if check.ID == id && check.ID != "" {
		//TODO: ID already exists, send 4xx error
		log.Printf("ID: %s already exists.", id)
		//TODO:
		return
	}

	log.Printf("Poll with ID: %s created.", id)

	//TODO: parse all poll fields

	ctx.String(http.StatusOK, "poll created: " + id)
}

func AddAllPollHandler(server *server.WebServer) {
	pollGroup := server.Router.Group("/api/v1/poll")

	pollGroup.POST("/create", createPollPOST)
}
