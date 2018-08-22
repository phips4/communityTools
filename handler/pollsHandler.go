package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/phips4/communityTools/app/db"
	"github.com/phips4/communityTools/app/logic"
	"github.com/phips4/communityTools/app/polls"
	"github.com/phips4/communityTools/server"
	"gopkg.in/mgo.v2"
	"net/http"
	"strconv"
)

var statusIllegalId = "id contains illegal character(s) or is too long. (a-zA-Z0-9_) 1-" + strconv.Itoa(logic.MaxStringLength)

const (
	statusInvalidPostParams = "invalid post parameter"
	statusPollNotExist      = "poll does not exist"
	statusErrorUpdating     = "error while updating poll"
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
func createPollPOST(ctx *gin.Context) (int, error) {
	id := ctx.PostForm("id")
	if !logic.ValidateID(id) {
		return http.StatusBadRequest, errors.New(statusIllegalId)
	}

	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	cookieCheck := ctx.PostForm("cookieCheck")
	multipleOptions := ctx.PostForm("multipleOptions")
	options := ctx.PostFormArray("options")
	del := ctx.PostForm("deleteIn")
	if !logic.ValidatePostParams(title, description, cookieCheck, multipleOptions, del, options) {
		return http.StatusBadRequest, errors.New(statusInvalidPostParams)
	}

	sess := db.GetPollSession()
	defer sess.Close()
	exist, err := sess.PollExists(id)
	if err != nil {
		if err == mgo.ErrNotFound {
			exist = false
		} else {
			return http.StatusInternalServerError, errors.New("error while fetching data from database")
		}
	}

	if exist {
		return http.StatusConflict, errors.New("poll ID already exists")
	}

	//64^DeleteIdLength should definitely be enough to avoid bruteforcing or guessing
	//64^7 = 4.398e+12
	rndStr, err := logic.GenerateRandomString(logic.DeleteIdLength)
	if err != nil {
		return http.StatusInternalServerError, errors.New("error while generating delete ID")
	}

	p := polls.NewPoll(id, title, description, cookieCheck, multipleOptions, rndStr, del, options)

	if err = sess.InsertPoll(p); err != nil {
		return http.StatusInternalServerError, errors.New("error saving data into database")
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "ok", "msg": "successfully created", "editToken": p.EditToken})
	return http.StatusOK, nil
}

/* +--------------------+
 * |        GET         |
 * +--------------------+
 * | • id               |
 * +--------------------+
 */
func getPollGET(ctx *gin.Context) (code int, err error) {
	var id string
	if code, err := getId(ctx, &id); err != nil {
		return code, err
	}

	sess := db.GetPollSession()
	defer sess.Close()

	poll, err := sess.GetPoll(id)
	if code, err := checkGetPoll(err); err != nil {
		return code, err
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "data": poll})
	return http.StatusOK, nil
}

/* +--------------------+
 * |        PUT         |
 * +--------------------+
 * | • id               |
 * | • cookieToken      |
 * | • option           |
 * +--------------------+
 */
func votePollPUT(ctx *gin.Context) (int, error) {
	var id string
	if code, err := getId(ctx, &id); err != nil {
		return code, err
	}

	sess := db.GetPollSession()
	defer sess.Close()
	poll, err := sess.GetPoll(id)

	if code, err := checkGetPoll(err); err != nil {
		return code, err
	}

	if poll.VotingStopped {
		return http.StatusBadRequest, errors.New("voting has already stopped")
	}

	cookieToken := ctx.PostForm("cookieToken")
	option := ctx.PostForm("option")
	//apply the new vote to the struct, so we can update it in the DB later
	if err := logic.ApplyVote(poll, ctx.ClientIP(), cookieToken, option); err != nil {
		if err == logic.ErrAlreadyVoted {
			return http.StatusBadRequest, errors.New("you have already voted")
		} else {
			return http.StatusBadRequest, errors.New("your given option might be wrong") //meh
		}
	}

	err = sess.UpdatePoll(id, poll)
	if err != nil {
		return http.StatusInternalServerError, errors.New(statusErrorUpdating)
	}

	ok(ctx, "successful voted for "+option+".")
	return http.StatusOK, nil
}

/* +--------------------+
 * |       PATCH        |
 * +--------------------+
 * | • id               |
 * | • editToken        |
 * +--------------------+
 */
func stopPollPATCH(ctx *gin.Context) (int, error) {
	var id string
	if code, err := getId(ctx, &id); err != nil {
		return code, err
	}

	sess := db.GetPollSession()
	defer sess.Close()

	p, err := sess.GetPoll(id)
	if code, err := checkGetPoll(err); err != nil {
		return code, err
	}

	edit := ctx.PostForm("editToken")
	if edit != p.EditToken {
		return http.StatusBadRequest, errors.New("invalid editToken")
	}

	p.VotingStopped = true
	err = sess.UpdatePoll(id, p)
	if err != nil {
		return http.StatusInternalServerError, errors.New(statusErrorUpdating)
	}

	ok(ctx, "voting stopped.")
	return http.StatusOK, nil
}

/* +--------------------+
 * |       DELETE       |
 * +--------------------+
 * | • id               |
 * | • editToken        |
 * | • sure             |
 * +--------------------+
 */
func deletePollDELETE(ctx *gin.Context) (int, error) {
	var id string
	if code, err := getId(ctx, &id); err != nil {
		return code, err
	}

	sure, err := strconv.ParseBool(ctx.PostForm("sure"))
	if err != nil {
		return http.StatusBadRequest, errors.New("can not convert post parameter 'sure' to bool")
	}
	if !sure {
		return http.StatusBadRequest, errors.New("deletion is not accepted by request")
	}

	sess := db.GetPollSession()
	defer sess.Close()
	p, err := sess.GetPoll(id)
	if code, err := checkGetPoll(err); err != nil {
		return code, err
	}

	edit := ctx.PostForm("editToken")
	if edit != p.EditToken {
		return http.StatusBadRequest, errors.New("invalid editToken")
	}

	err = sess.DeletePoll(id)
	if err != nil {
		return http.StatusInternalServerError, errors.New("error while deleting poll")
	}

	ok(ctx, "successfully deleted.")
	return http.StatusOK, nil
}

/**************************
 *  helper functions
 **************************/
 // adapter function for better error handling
func errorHandler(f func(ctx *gin.Context) (code int, err error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code, err := f(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"status": "error", "msg": err.Error()})
		}
	}
}

// returns true if id is valid
func getId(ctx *gin.Context, id *string) (int, error) {
	*id = ctx.Param("id")

	if !logic.ValidateID(*id) {
		return http.StatusBadRequest, errors.New(statusIllegalId)
	}

	return -1, nil
}

// returns true if no error occurred
func checkGetPoll(err error) (int, error) {
	if err != nil {
		if err == mgo.ErrNotFound {
			return http.StatusBadRequest, errors.New(statusPollNotExist)
		}

		return http.StatusInternalServerError, errors.New("error while searching ID from database")
	}

	return -1, nil
}

// register all poll endpoints
func AddAllPollHandler(server *server.WebServer) {
	pollGroup := server.Router.Group("/api/v1/poll")

	pollGroup.POST("/create", errorHandler(createPollPOST))
	pollGroup.PUT("/vote/:id", errorHandler(votePollPUT))
	pollGroup.GET("/get/:id", errorHandler(getPollGET))
	pollGroup.PATCH("/stop/:id", errorHandler(stopPollPATCH))
	pollGroup.DELETE("/delete/:id", errorHandler(deletePollDELETE))
}
