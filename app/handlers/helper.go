package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ok(ctx *gin.Context, msg string) {
	if msg == "" {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "msg": msg})
	}
}
