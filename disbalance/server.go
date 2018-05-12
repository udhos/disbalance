package main

import (
	"net/http"
	"sync"

	"gopkg.in/yaml.v2"
)

type rule struct {
	Name string
}

type config struct {
	basicAuthUser string
	basicAuthPass string
	rules         []rule
}

type server struct {
	cfg  config
	apis []string
	lock sync.RWMutex
}

func (s *server) auth(user, pass string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return user == s.cfg.basicAuthUser && pass == s.cfg.basicAuthPass
}

func (s *server) apiList() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.apis
}

func (s *server) ruleList() ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return yaml.Marshal(s.cfg.rules)
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
