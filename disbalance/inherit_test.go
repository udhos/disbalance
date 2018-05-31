package main

import (
	"testing"
)

func TestHost(t *testing.T) {
	testGetHost(t, "", "")
	testGetHost(t, ":", "")
	testGetHost(t, "a", "a")
	testGetHost(t, "a:b", "a")
	testGetHost(t, "a:", "a")
	testGetHost(t, "[::1]", "[::1]")
	testGetHost(t, "[::1]:2", "[::1]")
	testGetHost(t, "[::1]:", "[::1]")
	testGetHost(t, ":2", "")
}

func TestPort(t *testing.T) {
	testGetPort(t, "", "")
	testGetPort(t, ":", "")
	testGetPort(t, "a", "")
	testGetPort(t, "a:b", ":b")
	testGetPort(t, "a:", "")
	testGetPort(t, "[::1]", "")
	testGetPort(t, "[::1]:2", ":2")
	testGetPort(t, "[::1]:", "")
	testGetPort(t, ":2", ":2")
}

func testGetHost(t *testing.T, addr, expectedHost string) {
	h := getHost(addr)
	if h != expectedHost {
		t.Errorf("testGetHost: addr=[%s] expectedHost=[%s] gotHost=[%s]", addr, expectedHost, h)
	}
}

func testGetPort(t *testing.T, addr, expectedPort string) {
	p := getPort(addr)
	if p != expectedPort {
		t.Errorf("testGetPort: addr=[%s] expectedPort=[%s] gotPort=[%s]", addr, expectedPort, p)
	}
}
