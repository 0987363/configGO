package handlers

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/0987363/configGO/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var RootMux = gin.New()

func Init() {
	RootMux.Use(gin.Logger())

	work := viper.GetString("work")
	if work == "" {
		return
	}

	dirs, err := ioutil.ReadDir(work)
	if err != nil {
		log.Fatal("Read work failed:", err)
	}

	m := make(map[string]interface{})
	for _, p := range dirs {
		m[p.Name()] = ReadProject(filepath.Join(work, p.Name()))
	}

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
	for _, file := range files {
		k, v := ReadService(projectPath, file.Name())
		project[k] = v
//		log.Infof("Project: %s, Service:%s, data:%v", projectPath, k, v)
	}

	return project
}

func ReadService(dir, name string) (string, map[string]interface{}) {
	file := filepath.Join(dir, name)
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)
	if ext[0] == '.' {
		ext = ext[1:]
	}

	switch strings.ToLower(ext) {
	case "json":
		return name, ReadJsonService(file)
	case "toml":
		return name, ReadTomlService(file)
	default:
		return name, nil
	}
}

func ReadTomlService(file string) map[string]interface{} {
	return nil
}

func ReadJsonService(file string) map[string]interface{} {
	data, _ := ioutil.ReadFile(file)

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		log.Fatalf("Unmarshal service:%s failed:%v", file, err)
	}

	return m
}
