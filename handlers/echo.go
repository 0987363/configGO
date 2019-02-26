package handlers

import (
	"net/http"

	"github.com/0987363/configGO/middleware"
	"github.com/gin-gonic/gin"
)

func Echo(c *gin.Context) {
	v := middleware.GetConfig(c)
	c.JSON(http.StatusOK, v)
}
