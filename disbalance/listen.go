package main

import (
	"log"
	"net"
	"time"
)

func newInactiveTimer() *time.Timer {
	t := time.NewTimer(time.Second)
	if !t.Stop() {
		<-t.C // drain https://golang.org/pkg/time/#Timer.Reset
	}
	return t
}

func service_listen(ruleName, proto, listen string, enable chan bool, conn chan net.Conn) {
	log.Printf("listen: rule=%s proto=%s listen=%s starting", ruleName, proto, listen)

	listenRetry := time.Duration(3) * time.Second
	listenTimeout := time.Duration(3) * time.Second
	var tcpL *net.TCPListener
	listenTimer := newInactiveTimer()
	acceptTimer := newInactiveTimer()

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
			log.Printf("listen: rule=%s proto=%s listen=%s listener created: %v", ruleName, proto, listen, tcpL)
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
			log.Printf("listen: rule=%s proto=%s listen=%s new connection: %v", ruleName, proto, listen, newConn)
			conn <- newConn // send new conn to forwarder
		}
	}

	log.Printf("listen: rule=%s proto=%s listen=%s stopping", ruleName, proto, listen)
}
