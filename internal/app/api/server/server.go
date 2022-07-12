package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

type ApiServer interface {
	Start()
}
type (
	apiServer struct {
		srv *http.Server
	}
)

func NewServer(server *http.Server) ApiServer {
	return &apiServer{
		srv: server,
	}
}

func (api *apiServer) Start() {
	go func() {
		// service connections
		if err := api.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen server err,err:%s", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := api.srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server Shutdown,err:%s", err.Error())
	}
}
