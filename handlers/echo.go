package handlers

import (
	"net/http"

	"github.com/0987363/map-md5"
	"github.com/0987363/configGO/middleware"
	"github.com/gin-gonic/gin"
)

const HeaderEtag = "X-Druid-Etag"

func Echo(c *gin.Context) {
	v := middleware.GetConfig(c)
	if m, ok := v.(map[string]interface{}); ok {
		c.Header(HeaderEtag, mmd5.MapMd5(m))
		c.JSON(http.StatusOK, v)
		return
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
