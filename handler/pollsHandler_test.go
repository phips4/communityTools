package handler

import (
	"testing"
	"net/http/httptest"
	"github.com/phips4/communityTools/server"
	"src/github.com/stretchr/testify/assert"
	"net/http"
)

func TestAllPollEndpoints(t *testing.T) {
	// create new poll
	// vote
	// get
	// stop
	// get
	// delete
	// get
}

//TODO: VERY WIP VERY MUCH TODO
func TestGetPollGET(t *testing.T) {

	ws := server.New()
	ws.Router.ForwardedByClientIP = true

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/polls/get/AIWDAWDwAD", nil)
	ws.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "expected 200; got %v.", w.Code)
}