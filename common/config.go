package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
)

type Worker map[string]interface{}

func (w Worker) ParseProject(c chan *Application) {
	for project, v := range w {
		w, ok := v.(Worker)
		if !ok {
			continue
		}
		for service, v := range w {
			w, ok := v.(Worker)
			if !ok {
				continue
			}
			c <- &Application{
				Project: project,
				Service: service,
				Value:   w.ToString(),
				Op:      watcher.Create,
			}
		}
	}
}

func (w Worker) SearchProject(project string) Worker {
	if m, ok := w[project]; ok {
		if m, ok := m.(Worker); ok {
			return m
		}
	}
	return nil
}

func (w Worker) SearchService(project, service string) *Worker {
	m, ok := w[project]
	if !ok {
		return nil
	}
	w, ok = m.(Worker)
	if !ok {
		return nil
	}

	m, ok = w[service]
	if !ok {
		return nil
	}
	w, ok = m.(Worker)
	if !ok {
		return nil
	}

	return &w
}

func trimFileName(file string) string {
	return file[0 : len(file)-len(filepath.Ext(file))]
}

func splitPath(path string) (string, string) {
	dir, name := filepath.Split(path)
	project := filepath.Base(dir)
	return project, trimFileName(name)
}

func (w Worker) ToString() string {
	if w == nil {
		return ""
	}

	data, _ := json.Marshal(w)
	return string(data)
}

func (w Worker) Update(path string, c chan *Application) error {
	m := readService(path)
	if m == nil {
		return errors.New("Read service failed:" + path)
	}

	project, service := splitPath(path)
//	p := w.SearchProject(project)
	s := w.SearchService(project, service)
	if s == nil {
		return errors.New("Read project failed:" + path)
	}

	if reflect.DeepEqual(m, *s) {
		return nil
	}

	log.Info("Found diff:", m, *s)

	*s = m
	log.Info("Start send new app:", project, service)
	c <- &Application{
		Project: project,
		Service: service,
		Value:   m.ToString(),
		Op:      watcher.Write,
	}

	return nil
}

func NewWorker(work string) Worker {
	if work == "" {
		return nil
	}

	projects, err := ioutil.ReadDir(work)
	if err != nil {
		log.Fatal("Read work failed:", err)
	}

	w := make(Worker)
	for _, p := range projects {
		if p.Name()[0] == '.' {
			continue
		}
		if !p.IsDir() {
			continue
		}
		w[p.Name()] = readApplication(filepath.Join(work, p.Name()))
	}

	return w
}

func readApplication(projectPath string) Worker {
	services, err := ioutil.ReadDir(projectPath)
	if err != nil {
		log.Fatal("Read project failed:", err, projectPath)
	}
	if len(services) == 0 {
		return Worker{}
	}

	w := make(Worker)
	for _, file := range services {
		if file.Name()[0] == '.' {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext == "" {
			continue
		}

		if m := readService(filepath.Join(projectPath, file.Name())); m != nil {
			w[trimFileName(file.Name())] = m
		}
	}

	return w
}

func readService(path string) Worker {
	data, _ := ioutil.ReadFile(path)
	if len(data) == 0 {
		return nil
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return parseJson(data)
	case ".toml":
		return parseToml(data)
	default:
		return nil
	}
}

func parseToml(data []byte) Worker {
	return nil
}

func parseJson(data []byte) Worker {
	var m Worker
	if err := json.Unmarshal(data, &m); err != nil {
		log.Fatalf("Unmarshal service failed:%v", err)
	}

	return m
}
