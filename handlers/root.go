package handlers

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"reflect"

	"github.com/0987363/configGO/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var RootMux = gin.New()

func Init() {
	RootMux.Use(gin.Logger())

	work := viper.GetString("work")
	log.Debug("Config path:", work)

	dirs, err := ioutil.ReadDir(work)
	if err != nil {
		log.Fatal("Read path failed:", err)
	}
	log.Debug("Project count:", len(dirs))

	m := make(map[string]interface{})
	for _, p := range dirs {
		log.Debug("Project:", p.Name())
		m[p.Name()] = ReadProject(path.Join(work, p.Name()))
	}

	log.Info("Config:", m)

	BuildRouter(m)
}

func BuildRouter(m map[string]interface{}) {
	for p, s := range m {
		pMux := RootMux.Group("/" + p)

		services, ok := s.(map[string]interface{})
		if !ok {
			log.Fatal("Convert to services failed:", s)
		}
		if len(services) == 0 {
			continue
		}
		for f, v := range services {
			sMux := pMux.Group("/" + f)
			sMux.GET("/", middleware.Config(v, Echo))

			if d, ok := v.(map[string]interface{}); ok {
				BuildUrl(sMux, d)
			}
			log.Debug("Value type:", reflect.TypeOf(v))

		}
	}
}

func BuildUrl(fieldMux *gin.RouterGroup, fields map[string]interface{}) {
	for f, v := range fields {
		fMux := fieldMux.Group("/" + f)
		fMux.GET("/", middleware.Config(v, Echo))

		if d, ok := v.(map[string]interface{}); ok {
			BuildUrl(fMux, d)
		}
	}
}

func ReadProject(projectPath string) map[string]interface{} {
	files, err := ioutil.ReadDir(projectPath)
	if err != nil {
		log.Fatal("Read project failed:", err)
	}
	if len(files) == 0 {
		return map[string]interface{}{}
	}

	project := make(map[string]interface{})
	log.Infof("Project: %s, Service count:%d", projectPath, len(files))
	for _, file := range files {
		log.Debugf("Project: %s, Service:%s", projectPath, file.Name())
		project[file.Name()] = ReadService(path.Join(projectPath, file.Name()))
	}

	return project
}

func ReadService(file string) map[string]interface{} {
	data, _ := ioutil.ReadFile(file)

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		log.Fatalf("Unmarshal service:%s failed:%v", file, err)
	}

	log.Info("Service:", m)

	return m
}
