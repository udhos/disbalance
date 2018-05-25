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
	checks []checker
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

	for t, target := range f.rule.Targets {
		c := checker{
			done: make(chan struct{}),
		}
		f.checks = append(f.checks, c)
		go service_check(ruleName, proto, t, target, c, f.health)
	}

	healthTable := map[string]struct{}{}

LOOP:
	for {
		select {
		case <-f.done:
			break LOOP
		case h := <-f.health:
			log.Printf("forward: rule=%s target=%s healthy=%v", ruleName, h.target, h.status)
			if h.status {
				healthTable[h.target] = struct{}{}
			} else {
				delete(healthTable, h.target)
			}
		}
	}
	log.Printf("forward: rule=%s stopping", ruleName)

	for _, c := range f.checks {
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
			if !status {
				up++
				down = 0
				if up >= target.Check.Minimum {
					health <- targetHealth{targetName, true}
					status = true
				}
			}
		} else {
			if status {
				down++
				up = 0
				if down >= target.Check.Minimum {
					health <- targetHealth{targetName, false}
					status = false
				}
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
