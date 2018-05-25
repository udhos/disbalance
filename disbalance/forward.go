package main

import (
	"log"

	"github.com/udhos/disbalance/rule"
)

type forwarder struct {
	rule rule.Rule
	done chan struct{}
}

// beware: these methods hold an exclusive lock on the server
//         the lock is inherited from the caller

func (s *server) forwardEnable(ruleName string) {
	log.Printf("forwardEnable: rule=%s", ruleName)

	r, foundR := s.cfg.Rules[ruleName]
	if !foundR {
		log.Printf("forwardEnable: not found rule: rule=%s", ruleName)
		return
	}

	f, foundF := s.fwd[ruleName]
	if foundF {
		log.Printf("forwardEnable: unexpected existing forwarder: rule=%s", ruleName)
		return
	}

	f = forwarder{
		rule: *r.Clone(),
		done: make(chan struct{}),
	}

	s.fwd[ruleName] = f

	go service_forward(ruleName, f)
}

func (s *server) forwardDisable(ruleName string) {
	log.Printf("forwardDisable: rule=%s", ruleName)

	_, foundR := s.cfg.Rules[ruleName]
	if !foundR {
		log.Printf("forwardDisable: not found rule: rule=%s", ruleName)
		return
	}

	f, foundF := s.fwd[ruleName]
	if !foundF {
		log.Printf("forwardDisable: not found forwarder: rule=%s", ruleName)
		return
	}

	close(f.done) // request termination

	delete(s.fwd, ruleName)
}

func service_forward(ruleName string, f forwarder) {
	log.Printf("forwarder: rule=%s starting", ruleName)
LOOP:
	for {
		select {
		case <-f.done:
			break LOOP
		}
	}
	log.Printf("forwarder: rule=%s stopping", ruleName)
}
