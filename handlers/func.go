package handlers

import (
	"net/http"

	"github.com/0987363/configGO/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Echo(c *gin.Context) {
	v := middleware.GetConfig(c)
	log.Debug("Data:", v)
	c.JSON(http.StatusOK, v)
}
