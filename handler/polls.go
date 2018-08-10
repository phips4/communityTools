package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/server"
			"log"
	"net/http"
	"github.com/phips4/communityTools/app/db"
	"regexp"
	"gopkg.in/mgo.v2"
)

func createPollPOST(ctx *gin.Context) {
	sess := db.GetPollSession()
	defer sess.Close()

	id := ctx.DefaultPostForm("id", "")
	//TODO: post params verify

	ok, err := regexp.Match("^[a-zA-Z0-9_]*$", []byte(id))
	if err != nil {
		AbortWithError(ctx, http.StatusInternalServerError, "error while checking id")
		return
	}

	if !ok {
		AbortWithError(ctx, http.StatusBadRequest, "ID contains illegal character(s). (a-zA-Z0-9_)")
		return
	}

	exist, err := sess.PollExists(id)
	if err != nil {
		if err == mgo.ErrNotFound {
			exist = false
		} else {
			AbortWithError(ctx, http.StatusInternalServerError, "error while fetching data from database, ")
			return
		}
	}

	if exist {
		log.Println("poll id already exists")
		AbortWithError(ctx, http.StatusBadRequest, "Poll ID already exists")
		return
	}

	/*
	title := ctx.DefaultPostForm("title", "")
	description := ctx.DefaultPostForm("description", "")
	cookieCheck := ctx.DefaultPostForm("cookieCheck", "false")
	multipleOptions := ctx.DefaultPostForm("multipleOptions", "false")
	options := ctx.PostFormArray("options")
	*/

	ctx.JSON(http.StatusCreated, gin.H{"msg": "poll created"})

	//TODO: debug error -.-'

	//TODO: parse all poll fields
	//TODO: check all fields (400 response if not valid)
	//TODO: save poll in database
}

func AddAllPollHandler(server *server.WebServer) {
	pollGroup := server.Router.Group("/api/v1/poll")

	pollGroup.POST("/create", createPollPOST)
}
