package handler

import (
	"github.com/gin-gonic/gin"
)

func AbortWithError(ctx *gin.Context, code int, err string) {
	ctx.JSON(code, gin.H{"msg": err})
	ctx.Abort()
}
