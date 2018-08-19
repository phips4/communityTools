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

var statusIllegalId = "ID contains illegal character(s) or is too long. (a-zA-Z0-9_) 1-" + strconv.Itoa(logic.MaxStringLength)

const (
	statusInvalidPostParams = "invalid post parameter"
	statusPollNotExist      = "poll does not exist"
)

/* +--------------------+
 * |        POST        |
 * +--------------------+
 * | • id               |
 * | • title            |
 * | • description      |
 * | • cookieCheck      |
 * | • multipleOptions  |
 * | • options          |
 * | • deleteIn         |
 * +--------------------+
 */
func createPollPOST(ctx *gin.Context) {
	id := ctx.PostForm("id")
	if !logic.ValidateID(id) {
		AbortWithError(ctx, http.StatusBadRequest, statusIllegalId)
		return
	}

	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	cookieCheck := ctx.PostForm("cookieCheck")
	multipleOptions := ctx.PostForm("multipleOptions")
	options := ctx.PostFormArray("options")
	if !logic.ValidatePostParams(title, description, cookieCheck, multipleOptions, options) {
		AbortWithError(ctx, http.StatusBadRequest, statusInvalidPostParams)
		return
	}

	delete := ctx.PostForm("deleteIn")
	deleteIn, err := strconv.ParseInt(delete, 10, 64)
	if err != nil {
		AbortWithError(ctx, http.StatusBadRequest, "deleteIn is not a number.")
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

	//64^DeleteIdLength should definitely be enough to avoid bruteforcing or guessing
	//64^7 = 4.398e+12
	rndStr, err := logic.GenerateRandomString(logic.DeleteIdLength)
	if err != nil {
		AbortWithError(ctx, http.StatusInternalServerError, "error while generating delete ID.")
		return
	}

	// now, all parameters are valid
	p := polls.NewPoll(id, title, description, cookieCheck, multipleOptions, rndStr, int(deleteIn), options)

	if err = sess.SavePoll(p); err != nil {
		AbortWithError(ctx, http.StatusInternalServerError, "error saving data into database")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"msg": "poll successfully created"})
}

/* +--------------------+
 * |        GET         |
 * +--------------------+
 * | • id               |
 * +--------------------+
 */
func getPollGET(ctx *gin.Context) {
	id := ctx.Param("id")

	if !logic.ValidateID(id) {
		AbortWithError(ctx, http.StatusBadRequest, statusIllegalId)
		return
	}

	sess := db.GetPollSession()
	defer sess.Close()
	poll, err := sess.GetPoll(id)
	if err != nil {
		if err == mgo.ErrNotFound {
			AbortWithError(ctx, http.StatusBadRequest, statusPollNotExist)
			return
		}

		AbortWithError(ctx, http.StatusInternalServerError, "Error while searching ID from database.")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "ok", "data": poll})
}

/* +--------------------+
 * |        PUT         |
 * +--------------------+
 * | • id               |
 * | • cookieToken      |
 * +--------------------+
 */
func votePollPUT(ctx *gin.Context) {
	id := ctx.Param("id")
	//ip := ctx.Request.RemoteAddr
	//TODO: rewrite IP getter for production (headers, proxies, etc)
	ip := logic.GetIp(ctx.Request)
	cookieToken := ctx.PostForm("cookieToken")
	option := ctx.PostForm("option")

	if !logic.ValidateID(id) {
		AbortWithError(ctx, http.StatusBadRequest, statusIllegalId)
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

	if poll.VotingStopped {
		AbortWithError(ctx, http.StatusBadRequest, "voting has already stopped.")
		return
	}

	//apply the new vote to the struct, so we can update it in the DB later
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
	pollGroup.PUT("/vote/:id", votePollPUT)
	pollGroup.GET("/get/:id", getPollGET)
}
