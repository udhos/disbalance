package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/udhos/disbalance/rule"
)

func ruleUpdate(old, update rule.Rule) rule.Rule {
	r := rule.Rule{
		Name:     old.Name,
		Protocol: old.Protocol,
		Listener: old.Listener,
		Targets:  map[string]rule.Target{},
	}
	// copy from old
	for a, t := range old.Targets {
		r.Targets[a] = t
	}
	// copy from new
	if update.Protocol != "" {
		r.Protocol = update.Protocol
	}
	if update.Listener != "" {
		r.Listener = update.Listener
	}
	for a, t := range update.Targets {
		r.Targets[a] = t
	}
	return r
}

type config struct {
	BasicAuthUser string
	BasicAuthPass string
	Rules         map[string]rule.Rule
}

type server struct {
	cfg        config
	configPath string
	apis       []string
	lock       sync.RWMutex
}

// get a lock before calling unsafeSave
// will read from shared config
// will write into shared file
func unsafeSave(cfg *config, configPath string) {
	buf, errYaml := yaml.Marshal(cfg)
	if errYaml != nil {
		log.Printf("configSave: marshal: %s: %v", configPath, errYaml)
		return
	}
	if err := ioutil.WriteFile(configPath, buf, 0640); err != nil {
		log.Printf("configSave: %s: %v", configPath, err)
	}
	log.Printf("configSave: %s", configPath)
}

func (s *server) configSave() {
	s.lock.Lock() // exclusive lock: will write on shared file s.configPath
	defer s.lock.Unlock()

	unsafeSave(&s.cfg, s.configPath)
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

func (s *server) ruleList() []rule.Rule {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var rules []rule.Rule
	for _, r := range s.cfg.Rules {
		rules = append(rules, r)
	}

	return rules
}

func (s *server) ruleDump() ([]byte, error) {
	return yaml.Marshal(s.ruleList())
}

func (s *server) ruleGet(name string) (rule.Rule, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, r := range s.cfg.Rules {
		if r.Name == name {
			return r, nil
		}
	}
	return rule.Rule{}, fmt.Errorf("rule not found")
}

func (s *server) ruleDel(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.cfg.Rules, name)

	unsafeSave(&s.cfg, s.configPath)
}

func (s *server) rulePut(r rule.Rule) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.cfg.Rules[r.Name] = r // fully replace old rule, if any

	unsafeSave(&s.cfg, s.configPath)
}

func (s *server) rulePost(rules []rule.Rule) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, newRule := range rules {

		if oldRule, found := s.cfg.Rules[newRule.Name]; found {
			// new rule found
			// update old rule
			update := ruleUpdate(oldRule, newRule)
			s.cfg.Rules[newRule.Name] = update
			continue
		}

		// new rule not found
		// append new rule
		s.cfg.Rules[newRule.Name] = newRule
	}

	unsafeSave(&s.cfg, s.configPath)
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
