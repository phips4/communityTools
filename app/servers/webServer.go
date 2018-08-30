package servers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Webserver struct {
	Router *gin.Engine
	srv    *http.Server
	close  chan bool
}

//TODO:
func (w *Webserver) Listen(addr string) {
	w.srv = &http.Server{
		Addr:    addr,
		Handler: w.Router,
	}
	w.Router.ForwardedByClientIP = true

	go func() {
		if err := w.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("webserver (%s) listen: %s\n", addr, err)
		}
	}()
}

func (w *Webserver) Stop(ctx context.Context) error {
	return w.srv.Shutdown(ctx)
}
