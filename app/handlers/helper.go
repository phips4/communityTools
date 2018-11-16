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

// adapter function for better error handling
func errorHandler(f func(ctx *gin.Context) (code int, err error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code, err := f(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"status": "error", "msg": err.Error()})
		}
	}
}