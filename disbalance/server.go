package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"gopkg.in/yaml.v2"
)

type healthCheck struct {
	Interval int
	Timeout  int
	Address  string // if empty defaults to target address
}

type target struct {
	Address string
	Check   healthCheck
}

type rule struct {
	Name      string
	Protocol  string
	Listeners []string
	Targets   []target
}

type config struct {
	BasicAuthUser string
	BasicAuthPass string
	Rules         []rule
}

type server struct {
	cfg        config
	configPath string
	apis       []string
	lock       sync.RWMutex
}

func (s *server) configSave() {
	s.lock.Lock() // exclusive lock: will write on shared file s.configPath
	defer s.lock.Unlock()
	buf, errYaml := yaml.Marshal(s.cfg)
	if errYaml != nil {
		log.Printf("configSave: marshal: %s: %v", s.configPath, errYaml)
		return
	}
	if err := ioutil.WriteFile(s.configPath, buf, 0777); err != nil {
		log.Printf("configSave: %s: %v", s.configPath, err)
	}
	log.Printf("configSave: %s", s.configPath)
}

func (s *server) configLoad() {
	s.lock.Lock() // exclusive lock: will write on shared s.cfg
	defer s.lock.Unlock()
	buf, errRead := ioutil.ReadFile(s.configPath)
	if errRead != nil {
		log.Printf("configLoad: %s: %v", s.configPath, errRead)
	}
	if err := yaml.Unmarshal(buf, &s.cfg); err != nil {
		log.Printf("configLoad: unmarshal: %s: %v", s.configPath, err)
		return
	}
	log.Printf("configLoad: %s", s.configPath)
}

func (s *server) auth(user, pass string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return user == s.cfg.BasicAuthUser && pass == s.cfg.BasicAuthPass
}

func (s *server) apiList() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.apis
}

func (s *server) ruleList() []rule {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.cfg.Rules
}

func (s *server) ruleDump() ([]byte, error) {
	return yaml.Marshal(s.ruleList())
}

func (s *server) ruleGet(name string) (rule, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, r := range s.cfg.Rules {
		if r.Name == name {
			return r, nil
		}
	}
	return rule{}, fmt.Errorf("rule not found")
}

func auth(w http.ResponseWriter, r *http.Request, app *server) bool {

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	username, password, authOK := r.BasicAuth()
	if !authOK {
		http.Error(w, "Not authorized", 401)
		return false
	}

	if !app.auth(username, password) {
		http.Error(w, "Not authorized", 401)
		return false
	}

	return true
}
