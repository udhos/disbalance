package main

import (
	"log"
	"net"

	"github.com/udhos/disbalance/rule"
)

type targetHealth struct {
	target string
	status bool
}

type forwarder struct {
	rule   rule.Rule
	done   chan struct{}
	health chan targetHealth
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
		rule:   *r.Clone(),
		done:   make(chan struct{}),
		health: make(chan targetHealth),
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

	close(f.done) // request goroutine termination

	delete(s.fwd, ruleName)
}

func inheritPort(child, parent string) string {
	if getPort(child) == "" {
		child += getPort(parent) // force port into child
	}
	return child
}

func service_forward(ruleName string, f forwarder) {
	log.Printf("forward: rule=%s starting", ruleName)

	// spawn listener
	listenEnable := make(chan bool)   // send requests to listener
	listenConn := make(chan net.Conn) // receive connections from listener
	go service_listen(ruleName, f.rule.Protocol, f.rule.Listener, listenEnable, listenConn)

	// spawn health checkers
	var checks []checker
	for t, target := range f.rule.Targets {
		c := checker{
			done: make(chan struct{}),
		}
		checks = append(checks, c)

		t = inheritPort(t, f.rule.Listener)

		go service_check(ruleName, f.rule.Protocol, t, target, c, f.health)
	}

	healthTable := map[string]struct{}{}
	healthyTargets := 0
LOOP:
	for {
		select {
		case <-f.done:
			break LOOP
		case h := <-f.health:
			if h.status {
				healthTable[h.target] = struct{}{}
				healthyTargets++
				if healthyTargets == 1 {
					listenEnable <- true // enable listener
				}
			} else {
				delete(healthTable, h.target)
				healthyTargets--
				if healthyTargets == 0 {
					listenEnable <- false // disable listener
				}
			}
			log.Printf("forward: rule=%s target=%s status=%v healthyTargets=%d", ruleName, h.target, h.status, healthyTargets)
		case conn := <-listenConn:
			log.Printf("forward: rule=%s new connection from listener=%s: %v", ruleName, f.rule.Listener, conn)
		}
	}
	log.Printf("forward: rule=%s stopping", ruleName)

	// stop listener
	close(listenEnable)

	// stop health checkers
	for _, c := range checks {
		close(c.done) // request goroutine termination
	}
}
