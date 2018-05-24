package main

import (
	"log"
)

type forwarder struct {
}

// beware: these methods hold an exclusive lock on the server

func (s *server) forwardEnable(ruleName string) {
	log.Printf("forwardEnable: rule=%s", ruleName)
}

func (s *server) forwardDisable(ruleName string) {
	log.Printf("forwardDisable: rule=%s", ruleName)
}
