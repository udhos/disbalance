package main

import (
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
	cfg  config
	apis []string
	lock sync.RWMutex
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

func (s *server) ruleList() ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return yaml.Marshal(s.cfg.Rules)
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
