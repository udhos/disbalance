package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/udhos/disbalance/rule"
)

const (
	version = "0.0"
)

func main() {

	var app server
	app.cfg = config{
		Rules: map[string]rule.Rule{},
	}
	app.fwd = map[string]forwarder{}

	log.Printf("version %s runtime %s GOMAXPROCS=%d", version, runtime.Version(), runtime.GOMAXPROCS(0))

	runDirEnv := os.Getenv("DISBALANCE_RUN")
	runDir := runDirEnv
	if runDir == "" {
		runDir = "run/"
	} else if runDir[len(runDir)-1] != '/' {
		runDir += "/"
	}
	log.Printf("env DISBALANCE_RUN=[%s] - run directory: [%s]", runDirEnv, runDir)

	var controlAddress, consoleDir string
	var key, cert string

	flag.StringVar(&key, "key", runDir+"key.pem", "TLS key file")
	flag.StringVar(&cert, "cert", runDir+"cert.pem", "TLS cert file")
	flag.StringVar(&app.configPath, "config", runDir+"disbalance.conf", "config path")
	flag.StringVar(&consoleDir, "console", runDir+"console", "console directory")
	flag.StringVar(&controlAddress, "addr", ":8080", "control address")
	flag.Parse()

	app.configLoad()

	if app.cfg.BasicAuthUser == "" && app.cfg.BasicAuthPass == "" {
		log.Printf("setting empty basic auth credentials to admin:admin")
		app.cfg.BasicAuthUser = "admin"
		app.cfg.BasicAuthPass = "admin"
	}

	app.configSave()

	tls := true

	if !fileExists(key) {
		log.Printf("TLS key file not found: %s - disabling TLS", key)
		tls = false
	}

	if !fileExists(cert) {
		log.Printf("TLS cert file not found: %s - disabling TLS", cert)
		tls = false
	}

	// root handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("root handler - not found: %s", r.URL.Path)
		http.Error(w, "Not found", 404)
	})

	registerAPI(&app, "/api/", serveAPI)
	registerAPI(&app, "/api/rule/", serveAPIRule)
	registerAPI(&app, "/api/check/", serveAPICheck)
	registerAPI(&app, "/api/conn/", serveAPIConn)

	registerStatic(&app, "/console/", consoleDir)

	app.enableForwarders()

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
