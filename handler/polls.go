package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/app/db"
	"github.com/phips4/communityTools/app/logic"
	"github.com/phips4/communityTools/app/polls"
	"github.com/phips4/communityTools/server"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

func createPollPOST(ctx *gin.Context) {
	sess := db.GetPollSession()
	defer sess.Close()

	id := ctx.PostForm("id")
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	cookieCheck := ctx.PostForm("cookieCheck")
	multipleOptions := ctx.PostForm("multipleOptions")
	options := ctx.PostFormArray("options")

	log.Printf("id: %s, title: %s, desc: %s, cookie: %s, multiple: %s, options: %s.", id, title, description, cookieCheck, multipleOptions, options)

	if !logic.ValidateID(id) {
		AbortWithError(ctx, http.StatusBadRequest, "ID contains illegal character(s) or is too long. (a-zA-Z0-9_) 1-32")
		return
	}

	if !logic.ValidatePostParams(title, description, cookieCheck, multipleOptions, options) {
		AbortWithError(ctx, http.StatusBadRequest, "invalid post parameter")
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
		log.Println("poll ID already exists")
		AbortWithError(ctx, http.StatusBadRequest, "Poll ID already exists")
		return
	}
	// now, all parameters are valid
	p := polls.NewPoll(id, title, description, cookieCheck, multipleOptions, options)

	if err = sess.SavePoll(p); err != nil {
		AbortWithError(ctx, http.StatusInternalServerError, "error saving data into database")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"msg": "poll successfully created"})
}

func AddAllPollHandler(server *server.WebServer) {
	pollGroup := server.Router.Group("/api/v1/poll")

	pollGroup.POST("/create", createPollPOST)
}
