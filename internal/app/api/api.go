package api

import (
	"net/http"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/api/rest"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/api/router"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/api/server"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
)

func NewApiServer(config *config.App) server.ApiServer {
	r := router.NewRouter()
	r.GET("/version", rest.GetVersion)
	srv := &http.Server{
		Addr:    config.Addr,
		Handler: r,
	}
	return server.NewServer(srv)

}
