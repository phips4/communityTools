package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/app/db"
	"github.com/phips4/communityTools/app/logic"
	"github.com/phips4/communityTools/app/polls"
	"github.com/phips4/communityTools/server"
	"gopkg.in/mgo.v2"
	"net/http"
	"strconv"
)

func createPollPOST(ctx *gin.Context) {
	id := ctx.PostForm("id")
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	cookieCheck := ctx.PostForm("cookieCheck")
	multipleOptions := ctx.PostForm("multipleOptions")
	options := ctx.PostFormArray("options")

	if !logic.ValidateID(id) {
		AbortWithError(ctx, http.StatusBadRequest, "ID contains illegal character(s) or is too long. (a-zA-Z0-9_) 1-" + strconv.Itoa(logic.MaxStringLength))
		return
	}

	if !logic.ValidatePostParams(title, description, cookieCheck, multipleOptions, options) {
		AbortWithError(ctx, http.StatusBadRequest, "invalid post parameter")
		return
	}

	sess := db.GetPollSession()
	defer sess.Close()

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

func getPollGET(ctx *gin.Context) {
	id := ctx.Param("id")

	if !logic.ValidateID(id) {
		AbortWithError(ctx, http.StatusBadRequest, "ID contains illegal character(s) or is too long. (a-zA-Z0-9_) 1-" + strconv.Itoa(logic.MaxStringLength))
		return
	}

	sess := db.GetPollSession()
	defer sess.Close()

	poll, err := sess.GetPoll(id)
	if err != nil {
		if err == mgo.ErrNotFound {
			AbortWithError(ctx, http.StatusBadRequest, "Poll does not exist")
			return
		}

		AbortWithError(ctx, http.StatusInternalServerError, "Error while searching ID from database.")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "ok", "data": poll})
}

func votePollPUT(ctx *gin.Context) {
	id := ctx.Param("id")
	//ip := ctx.Request.RemoteAddr
	//TODO: rewrite IP getter for produktion (headers, proxies, etc)
	ip := logic.GetIp(ctx.Request)
	cookieToken := ctx.PostForm("cookieToken")
	option := ctx.PostForm("option")

	if !logic.ValidateID(id) {
		AbortWithError(ctx, http.StatusBadRequest, "ID contains illegal character(s) or is too long. (a-zA-Z0-9_) 1-" + strconv.Itoa(logic.MaxStringLength))
		return
	}

	sess := db.GetPollSession()
	defer sess.Close()

	poll, err := sess.GetPoll(id)
	if err != nil {
		if err == mgo.ErrNotFound {
			AbortWithError(ctx, http.StatusNotFound, "Id not found")
			return
		}

		AbortWithError(ctx, http.StatusInternalServerError, "error while getting poll from database.")
		return
	}

	if !logic.ApplyVote(poll, ip, cookieToken, option) {
		AbortWithError(ctx, http.StatusBadRequest, "you already have voted or your option is bad.")
		return
	}

	err = sess.UpdatePoll(id, poll)
	if err != nil {
		AbortWithError(ctx, http.StatusInternalServerError, "error while updating poll.")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "successfully voted for " + option})
}

func AddAllPollHandler(server *server.WebServer) {
	pollGroup := server.Router.Group("/api/v1/poll")

	pollGroup.POST("/create", createPollPOST)
	pollGroup.GET("/get/:id", getPollGET)
	pollGroup.PUT("/vote/:id", votePollPUT)
}
