package middleware

import (
	"time"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/constant"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "HEAD", "DELETE", "PATCH"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type",
			constant.HeaderAuthorization, constant.HeaderXForwardedFor, "X-Real-Ip",
			"X-Appengine-Remote-Addr", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{constant.HeaderPagination, constant.HeaderContentDisposition},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))
}
