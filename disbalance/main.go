package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
)

const (
	version = "0.0"
)

func main() {

	var app server
	app.cfg.BasicAuthUser = "admin"
	app.cfg.BasicAuthPass = "admin"
	app.cfg.Rules = []rule{{Name: "rule0", Protocol: "tcp"}, {Name: "rule1", Protocol: "tcp"}}

	log.Printf("version %s runtime %s GOMAXPROCS=%d", version, runtime.Version(), runtime.GOMAXPROCS(0))

	var configPath, controlAddress, consoleDir string
	var key, cert string

	flag.StringVar(&key, "key", "key.pem", "TLS key file")
	flag.StringVar(&cert, "cert", "cert.pem", "TLS cert file")
	flag.StringVar(&configPath, "config", "run/disbalance.conf", "config path")
	flag.StringVar(&consoleDir, "console", "run/console", "console directory")
	flag.StringVar(&controlAddress, "addr", ":8080", "control address")
	flag.Parse()

	tls := true

	if !fileExists(key) {
		log.Printf("TLS key file not found: %s - disabling TLS", key)
		tls = false
	}

	if !fileExists(cert) {
		log.Printf("TLS cert file not found: %s - disabling TLS", cert)
		tls = false
	}

	registerApi(&app, "/api/", serveApi)
	registerApi(&app, "/api/rule/", serveApiRule)

	registerStatic(&app, "/console/", consoleDir)

	log.Printf("api credentials: user=%s pass=%s", app.cfg.BasicAuthUser, app.cfg.BasicAuthPass)

	if tls {
		log.Printf("serving HTTPS on TCP %s", controlAddress)
		if err := http.ListenAndServeTLS(controlAddress, cert, key, nil); err != nil {
			log.Fatalf("ListenAndServeTLS: %s: %v", controlAddress, err)
		}
		return
	}

	log.Printf("serving HTTP on TCP %s", controlAddress)
	if err := http.ListenAndServe(controlAddress, nil); err != nil {
		log.Fatalf("ListenAndServe: %s: %v", controlAddress, err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
