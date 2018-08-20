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
	statusErrorUpdateing = "error while updating poll."
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
	var id string
	if ok := getId(ctx, &id); !ok {
		return
	}

	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	cookieCheck := ctx.PostForm("cookieCheck")
	multipleOptions := ctx.PostForm("multipleOptions")
	options := ctx.PostFormArray("options")
	del := ctx.PostForm("deleteIn")
	if !logic.ValidatePostParams(title, description, cookieCheck, multipleOptions, del, options) {
		AbortWithError(ctx, http.StatusBadRequest, statusInvalidPostParams)
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
		AbortWithError(ctx, http.StatusConflict, "Poll ID already exists")
		return
	}

	//64^DeleteIdLength should definitely be enough to avoid bruteforcing or guessing
	//64^7 = 4.398e+12
	rndStr, err := logic.GenerateRandomString(logic.DeleteIdLength)
	if err != nil {
		AbortWithError(ctx, http.StatusInternalServerError, "error while generating delete ID.")
		return
	}

	p := polls.NewPoll(id, title, description, cookieCheck, multipleOptions, rndStr, del, options)

	if err = sess.InsertPoll(p); err != nil {
		AbortWithError(ctx, http.StatusInternalServerError, "error saving data into database")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "ok", "msg": "successfully created", "editToken": p.EditToken})
}

/* +--------------------+
 * |        GET         |
 * +--------------------+
 * | • id               |
 * +--------------------+
 */
func getPollGET(ctx *gin.Context) {
	var id string
	if ok := getId(ctx, &id); !ok {
		return
	}

	sess := db.GetPollSession()
	defer sess.Close()
	poll, err := sess.GetPoll(id)
	if !checkGetPoll(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "data": poll})
}

/* +--------------------+
 * |        PUT         |
 * +--------------------+
 * | • id               |
 * | • cookieToken      |
 * | • option           |
 * +--------------------+
 */
func votePollPUT(ctx *gin.Context) {
	var id string
	if ok := getId(ctx, &id); !ok {
		return
	}

	sess := db.GetPollSession()
	defer sess.Close()
	poll, err := sess.GetPoll(id)
	checkGetPoll(ctx, err)

	if poll.VotingStopped {
		AbortWithError(ctx, http.StatusBadRequest, "voting has already stopped.")
		return
	}

	cookieToken := ctx.PostForm("cookieToken")
	option := ctx.PostForm("option")
	//apply the new vote to the struct, so we can update it in the DB later
	if err := logic.ApplyVote(poll, ctx.ClientIP(), cookieToken, option); err != nil {
		if err == logic.ErrAlreadyVoted {
			AbortWithError(ctx, http.StatusBadRequest, "you have already voted")
		} else {
			AbortWithError(ctx, http.StatusBadRequest, "your given option might me wrong") //meh
		}
		return
	}

	err = sess.UpdatePoll(id, poll)
	if err != nil {
		AbortWithError(ctx, http.StatusInternalServerError, statusErrorUpdateing)
		return
	}

	ok(ctx, "successful voted for " + option + ".")
}

/* +--------------------+
 * |      PATCH         |
 * +--------------------+
 * | • id               |
 * | • editToken        |
 * +--------------------+
 */
func stopPollPATCH(ctx *gin.Context) {
	var id string
	if ok := getId(ctx, &id); !ok {
		return
	}

	sess := db.GetPollSession()
	defer sess.Close()

	p, err := sess.GetPoll(id)
	if !checkGetPoll(ctx, err) {
		return
	}

	edit := ctx.PostForm("editToken")
	if edit != p.EditToken {
		AbortWithError(ctx, http.StatusBadRequest, "invalid editToken")
		return
	}

	p.VotingStopped = true
	err = sess.UpdatePoll(id, p)
	if err != nil {
		AbortWithError(ctx, http.StatusInternalServerError, statusErrorUpdateing)
		return
	}

	ok(ctx, "voting stopped.")
}

/* helper functions */

// returns true if id is valid
func getId(ctx *gin.Context, id *string) bool {
	*id = ctx.Param("id")

	if !logic.ValidateID(*id) {
		AbortWithError(ctx, http.StatusBadRequest, statusIllegalId)
		return false
	}
	return true
}

// returns true if no error occurred
func checkGetPoll(ctx *gin.Context, err error) bool {
	if err != nil {
		if err == mgo.ErrNotFound {
			AbortWithError(ctx, http.StatusBadRequest, statusPollNotExist)
			return false
		}

		AbortWithError(ctx, http.StatusInternalServerError, "Error while searching ID from database.")
		return false
	}

	return true
}

func AddAllPollHandler(server *server.WebServer) {
	pollGroup := server.Router.Group("/api/v1/poll")

	pollGroup.POST("/create", createPollPOST)
	pollGroup.PUT("/vote/:id", votePollPUT)
	pollGroup.GET("/get/:id", getPollGET)
	pollGroup.PATCH("/stop/:id", stopPollPATCH)
}
