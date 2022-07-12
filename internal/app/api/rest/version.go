package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const version = "0.0.1"

type Version struct {
	Version string
}

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, Version{Version: version})
}
