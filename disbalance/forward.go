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

func service_listen(ruleName, proto, listen string, enable chan bool, conn chan net.Conn) {
	log.Printf("listen: rule=%s proto=%s listen=%s starting", ruleName, proto, listen)

	listenTimer := time.NewTimer(0)
	listenTimer.Stop() // create inactive timer
	listenRetry := time.Duration(3) * time.Second
	listenTimeout := time.Duration(3) * time.Second
	var tcpL *net.TCPListener
	acceptTimer := time.NewTimer(0)
	listenTimer.Stop() // create inactive timer

	stop := func() {
		acceptTimer.Stop()
		listenTimer.Stop()
		if tcpL != nil {
			tcpL.Close() // shutdown listener
			tcpL = nil
		}
	}

LOOP:
	for {
		select {
		case e, ok := <-enable:
			log.Printf("listen: rule=%s proto=%s listen=%s enable=%v", ruleName, proto, listen, e)
			if !ok {
				stop()
				break LOOP
			}
			if e {
				// request for start
				listenTimer.Reset(0) // schedule new listener
			} else {
				// request for stop
				stop()
			}
		case <-listenTimer.C:
			ln, errListen := net.Listen(proto, listen)
			if errListen != nil {
				log.Printf("listen: rule=%s proto=%s listen=%s error=%v", ruleName, proto, listen, errListen)
				listenTimer.Reset(listenRetry) // reschedule listen
				continue LOOP
			}
			var isTcp bool
			tcpL, isTcp = ln.(*net.TCPListener)
			if !isTcp {
				log.Printf("listen: rule=%s proto=%s listen=%s not tcp listener: %v", ruleName, proto, listen, ln)
				ln.Close()
				listenTimer.Reset(listenRetry) // reschedule listen
				continue LOOP
			}
			log.Printf("listen: rule=%s proto=%s listen=%s listener created", ruleName, proto, listen)
			acceptTimer.Reset(0) // schedule accept
		case <-acceptTimer.C:
			deadline := time.Now().Add(listenTimeout)
			if errDeadline := tcpL.SetDeadline(time.Now().Add(listenTimeout)); errDeadline != nil {
				log.Printf("listen: rule=%s proto=%s listen=%s deadline=%v error=%v", ruleName, proto, listen, deadline, errDeadline)
				tcpL.Close()
				listenTimer.Reset(listenRetry) // reschedule listen
				continue LOOP
			}
			newConn, errAccept := tcpL.Accept()
			acceptTimer.Reset(0) // reschedule accept
			if errAccept != nil {
				log.Printf("listen: rule=%s proto=%s listen=%s accept=%v", ruleName, proto, listen, errAccept)
				continue LOOP
			}
			log.Printf("listen: rule=%s proto=%s listen=%s new connection", ruleName, proto, listen)
			conn <- newConn // send new conn to forwarder
		}
	}

	log.Printf("listen: rule=%s proto=%s listen=%s stopping", ruleName, proto, listen)
}

func service_check(ruleName, proto, targetName string, target rule.Target, chk checker, health chan targetHealth) {
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
		conn, err := net.DialTimeout(proto, checkAddress, timeout)
		log.Printf("check: rule=%s target=%s check=%s err=%v up=%d down=%d status=%v", ruleName, targetName, checkAddress, err, up, down, status)
		if err == nil {
			conn.Close()
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
		case <-chk.done:
			break LOOP
		case <-ticker.C:
		}
	}

	ticker.Stop()

	log.Printf("check: rule=%s target=%s stopping", ruleName, targetName)
}
