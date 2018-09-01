package servers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type DefaultServer struct {
	Router *gin.Engine
	srv    *http.Server
	name   string
}

func NewDefaultServer() *DefaultServer {
	return &DefaultServer{
		gin.Default(),
		nil,
		"defaultServer",
	}
}

func (server *DefaultServer) Listen(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	server.srv = &http.Server{
		Addr:    addr,
		Handler: server.Router,
	}
	server.Router.ForwardedByClientIP = true

	log.Printf("%s server listening on :%d", server.name, port)
	if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("%s (%s) listen: %s\n", server.name, addr, err)
	}
}

func (server *DefaultServer) Stop(ctx context.Context) error {
	return server.srv.Shutdown(ctx)
}
