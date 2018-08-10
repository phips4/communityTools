package server

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func mgoSessionMiddleware(mgoSession *mgo.Session) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//session := mgoSession.Clone()
		//defer session.Close()

		//ctx.Set("mgo", session)
	}
}
