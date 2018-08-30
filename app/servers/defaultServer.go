package servers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type DefaultServer struct {
	Router *gin.Engine
	srv    *http.Server
}

func New() *DefaultServer {
	return &DefaultServer{
		gin.Default(),
		nil,
	}
}

func (server *DefaultServer) Listen(addr string) {
	server.srv = &http.Server{
		Addr:    addr,
		Handler: server.Router,
	}
	server.Router.ForwardedByClientIP = true

	if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("webserver (%s) listen: %s\n", addr, err)
	}
}

func (server *DefaultServer) Stop(ctx context.Context) error {
	return server.srv.Shutdown(ctx)
}
