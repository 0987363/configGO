package common

import (
	"encoding/json"
	"errors"
	"reflect"

	//	"errors"
	"io/ioutil"
	"path/filepath"

	//	"reflect"
	"strings"

	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	Map map[string]interface{}
	Ch  chan *Service
}

type Service struct {
	Project string
	Service string
	Map     map[string]interface{}
	Value   string
	Op      watcher.Op
}
type Project struct {
	Project string
	Map     map[string]interface{}
}

func (p *Project) RemoveService(key string) {
	p.Map[key] = nil
}

func (p *Project) GetService(key string) *Service {
	m, ok := p.Map[key]
	if !ok {
		m2 := make(map[string]interface{})
		p.Map[key] = m2
	}
	m2, ok := m.(map[string]interface{})
	if !ok {
		m2 = make(map[string]interface{})
		p.Map[key] = m2
	}
	return &Service{
		Project: p.Project,
		Service: key,
		Map:     m2,
	}
}

func (w *Worker) GetProject(key string) *Project {
	m, ok := w.Map[key]
	if !ok {
		m2 := make(map[string]interface{})
		w.Map[key] = m2
	}
	m2, ok := m.(map[string]interface{})
	if !ok {
		m2 = make(map[string]interface{})
		w.Map[key] = m2
	}
	return &Project{
		Project: key,
		Map:     m2,
	}
	return nil
}

func (w *Worker) GetService(project, service string) *Service {
	p := w.GetProject(project)
	s := p.GetService(service)
	if len(s.Map) == 0 {
		return nil
	}
	return s
}

func (w *Worker) Send(s *Service) {
	if len(s.Map) == 0 {
		return
	}

	s.Op = watcher.Write
	s.ToString()
	w.Ch <- s
}

func (w *Worker) ParseProject() {
	for project, v := range w.Map {
		for service, _ := range v.(map[string]interface{}) {
			s := w.GetService(project, service)
			w.Send(s)
		}
	}
}

func (w *Worker) SearchService(s *Service) *Service {
	p := w.GetProject(s.Project)
	return p.GetService(s.Service)
}

func (w *Worker) RemoveService(s *Service) {
	p := w.GetProject(s.Project)
	p.RemoveService(s.Service)
}

func (w *Worker) Remove(path string) error {
	s := NewService(path)

	w.RemoveService(s)
	s.Op = watcher.Remove
	w.Ch <- s

	return nil
}

func (w Worker) Update(path string) error {
	s := NewService(path)
	s, err := s.Load(path)
	if err != nil {
		return errors.New("Load service failed:" + path)
	}
	if len(s.Map) == 0 {
		return nil
	}

	old := w.SearchService(s)
	if reflect.DeepEqual(old.Map, s.Map) {
		return nil
	}

	w.AddService(s)
	w.Send(s)

	return nil
}

func NewWorker(dir string, c chan *Service) *Worker {
	if dir == "" {
		return nil
	}

	projects, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("Read work dir failed:", err)
	}

	w := &Worker{
		Map: make(map[string]interface{}),
		Ch:  c,
	}
	for _, p := range projects {
		if p.Name()[0] == '.' {
			continue
		}
		if !p.IsDir() {
			continue
		}
		prj := &Project{
			Project: p.Name(),
		}
		prj.LoadProject(w.AddService, filepath.Join(dir, p.Name()))
	}

	return w
}

func (w *Worker) AddService(s *Service) {
	m, ok := w.Map[s.Project]
	if !ok {
		m2 := make(map[string]interface{})
		w.Map[s.Project] = m2
	}

	m2, ok := m.(map[string]interface{})
	if !ok {
		m2 = make(map[string]interface{})
		w.Map[s.Project] = m2
	}

	m2[s.Service] = s.Map
}

func (n *Project) LoadProject(f func(*Service), path string) {
	services, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("Read project failed:", err, path)
	}
	if len(services) == 0 {
		log.Warning("Could not found service:", path)
		return
	}

	for _, file := range services {
		if file.Name()[0] == '.' {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext == "" {
			continue
		}

		s := &Service{
			Project: n.Project,
			Service: trimFileName(file.Name()),
		}
		s, err := s.Load(filepath.Join(path, file.Name()))
		if err != nil {
			log.Error("Load new service failed:", err)
			continue
		}
		if len(s.Map) != 0 {
			f(s)
		}
	}
}

func LoadService(path string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return parseJson(data), nil
	default:
		return nil, nil
	}
}

func NewService(path string) *Service {
	dir, name := filepath.Split(path)
	return &Service{
		Project: filepath.Base(dir),
		Service: trimFileName(name),
	}
}

func (s *Service) ToString() {
	data, _ := json.Marshal(s.Map)
	s.Value = string(data)
}

func (s *Service) Load(path string) (*Service, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return s, nil
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		s.Map = parseJson(data)
		return s, nil
	default:
		return s, nil
	}
}

func parseToml(data []byte) map[string]interface{} {
	return nil
}

func parseJson(data []byte) map[string]interface{} {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		log.Fatalf("Unmarshal service failed:%v", err)
	}

	return m
}

func trimFileName(file string) string {
	return file[0 : len(file)-len(filepath.Ext(file))]
}

func splitPath(path string) (string, string) {
	dir, name := filepath.Split(path)
	project := filepath.Base(dir)
	return project, trimFileName(name)
}
