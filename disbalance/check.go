package main

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/udhos/disbalance/rule"
)

type checker struct {
	done chan struct{}
}

// host:port => :port
func getPort(addr string) string {
	bracket := strings.LastIndexByte(addr, ']')
	addr = addr[bracket+1:]
	colon := strings.LastIndexByte(addr, ':')
	if colon < 0 {
		return ""
	}
	port := addr[colon:]
	if port == ":" {
		return ""
	}
	return port
}

func service_check(ruleName, proto, targetName string, target rule.Target, chk checker, health chan targetHealth) {
	log.Printf("check: rule=%s target=%s starting", ruleName, targetName)

	checkAddress := target.Check.Address

	// check inherits address from target
	if checkAddress == "" {
		checkAddress = targetName
	}

	checkAddress = inheritPort(checkAddress, targetName)

	log.Printf("check: rule=%s target=%s check=%s", ruleName, targetName, checkAddress)

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
