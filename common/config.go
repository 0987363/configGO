package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func ReadWork() map[string]interface{} {
	work := viper.GetString("work")
	if work == "" {
		return nil
	}

	dirs, err := ioutil.ReadDir(work)
	if err != nil {
		log.Fatal("Read work failed:", err)
	}

	m := make(map[string]interface{})
	for _, p := range dirs {
		if !p.IsDir() {
			continue
		}
		m[p.Name()] = readProject(filepath.Join(work, p.Name()))
	}

	return m
}

func readProject(projectPath string) map[string]interface{} {
	files, err := ioutil.ReadDir(projectPath)
	if err != nil {
		log.Fatal("Read project failed:", err, projectPath)
	}
	if len(files) == 0 {
		return map[string]interface{}{}
	}

	project := make(map[string]interface{})
	for _, file := range files {
		k, v := readService(projectPath, file.Name())
		if k == "" {
			continue
		}
		project[k] = v
	}

	return project
}

func readService(dir, name string) (string, map[string]interface{}) {
	ext := filepath.Ext(name)
	if len(ext) > 3 && ext[0] == '.' {
		ext = ext[1:]
	}
	file := filepath.Join(dir, name)
	name = strings.TrimSuffix(name, ext)

	switch strings.ToLower(ext) {
	case "json":
		return name, readJsonService(file)
	case "toml":
		return name, readTomlService(file)
	default:
		return name, nil
	}
}

func readTomlService(file string) map[string]interface{} {
	return nil
}

func readJsonService(file string) map[string]interface{} {
	data, _ := ioutil.ReadFile(file)

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		log.Fatalf("Unmarshal service:%s failed:%v", file, err)
	}

	return m
}