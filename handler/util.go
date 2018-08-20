package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AbortWithError(ctx *gin.Context, code int, err string) {
	ctx.JSON(code, gin.H{"status": err})
	ctx.Abort()
}

func ok(ctx *gin.Context, msg string) {
	if msg == "" {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "msg": msg})
	}
}