package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AbortWithError(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err})
	ctx.Abort()
}
