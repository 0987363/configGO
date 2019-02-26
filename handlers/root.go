package handlers

import (
	"github.com/0987363/configGO/common"
	"github.com/0987363/configGO/middleware"
	"github.com/gin-gonic/gin"
)

var RootMux = gin.New()

func Init() {
	RootMux.Use(gin.Recovery())
	RootMux.Use(middleware.RequestID())
	RootMux.Use(middleware.Logger())

	buildRouter(&RootMux.RouterGroup, common.ReadWork())
}

func buildRouter(mux *gin.RouterGroup, work map[string]interface{}) {
	for k, v := range work {
		fieldMux := mux.Group("/" + k)
		fieldMux.GET("/", middleware.Config(v, Echo))

		if d, ok := v.(map[string]interface{}); ok {
			buildRouter(fieldMux, d)
		}
	}
}

