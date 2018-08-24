package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type WebServer struct {
	Router *gin.Engine
	srv    *http.Server
}

func New() *WebServer {
	return &WebServer{
		gin.Default(),
		nil,
	}
}

func (server *WebServer) Listen(addr string) {
	server.srv = &http.Server{
		Addr:    addr,
		Handler: server.Router,
	}
	server.Router.ForwardedByClientIP = true

	go func() {
		if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("webserver (%s) listen: %s\n", addr, err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("good bye!")
}

func (server *WebServer) Stop(ctx context.Context) error {
	return server.srv.Shutdown(ctx)
}
