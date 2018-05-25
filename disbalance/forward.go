package main

import (
	"log"
	"net"
	"time"

	"github.com/udhos/disbalance/rule"
)

type checker struct {
	done chan struct{}
}

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

	go service_forward(ruleName, r.Protocol, f)
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

func service_forward(ruleName, proto string, f forwarder) {
	log.Printf("forward: rule=%s starting", ruleName)

	var checks []checker

	for t, target := range f.rule.Targets {
		c := checker{
			done: make(chan struct{}),
		}
		checks = append(checks, c)
		go service_check(ruleName, proto, t, target, c, f.health)
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
			} else {
				delete(healthTable, h.target)
				healthyTargets--
			}
			log.Printf("forward: rule=%s target=%s status=%v healthyTargets=%d", ruleName, h.target, h.status, healthyTargets)
		}
	}
	log.Printf("forward: rule=%s stopping", ruleName)

	for _, c := range checks {
		close(c.done) // request goroutine termination
	}
}

func service_check(ruleName, proto, targetName string, target rule.Target, c checker, health chan targetHealth) {
	log.Printf("check: rule=%s target=%s starting", ruleName, targetName)

	checkAddress := target.Check.Address
	if checkAddress == "" {
		checkAddress = targetName
	}

	timeout := time.Duration(target.Check.Timeout) * time.Second

	ticker := time.NewTicker(time.Duration(target.Check.Interval) * time.Second)

	var status bool
	var up, down int
LOOP:
	for {
		_, err := net.DialTimeout(proto, checkAddress, timeout)
		log.Printf("check: rule=%s target=%s check=%s err=%v up=%d down=%d status=%v", ruleName, targetName, checkAddress, err, up, down, status)
		if err == nil {
			down = 0
			up++
			if !status && up >= target.Check.Minimum {
				status = true
				health <- targetHealth{targetName, true}
			}
		} else {
			up = 0
			down++
			if status && down >= target.Check.Minimum {
				status = false
				health <- targetHealth{targetName, false}
			}
		}

		select {
		case <-c.done:
			break LOOP
		case <-ticker.C:
		}
	}

	ticker.Stop()

	log.Printf("check: rule=%s target=%s stopping", ruleName, targetName)
}
