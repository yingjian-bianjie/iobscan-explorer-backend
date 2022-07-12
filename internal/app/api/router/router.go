package router

import (
	"log"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/api/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Recovery())

	// log
	r.Use(gin.Logger())
	log.SetOutput(gin.DefaultWriter)

	middleware.SetCors(r)
	return r
}
