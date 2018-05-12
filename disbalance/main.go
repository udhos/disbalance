package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

const (
	version = "0.0"
)

type config struct {
}

type server struct {
	cfg config
}

func main() {

	var app server

	log.Printf("version %s runtime %s", version, runtime.Version())

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

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) { serveApi(w, r, &app) })

	registerStatic(&app, "/console/", consoleDir)

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

func serveApi(w http.ResponseWriter, r *http.Request, app *server) {
	log.Printf("serveApi: url=%s from=%s", r.URL.Path, r.RemoteAddr)

	io.WriteString(w,
		`<!DOCTYPE html>
<html lang="en-US">

<title>
<head>disbalance api</head>
</title>

<body>
<p>welcome to the api</p>
<a href="/console">console</a>
</body>

</html>
`)

}

type httpHandler struct {
	app *server
	h   http.Handler
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("staticHandler.ServeHTTP url=%s from=%s", r.URL.Path, r.RemoteAddr)
	h.h.ServeHTTP(w, r)
}

func registerStatic(app *server, path, dir string) {
	log.Printf("mapping www path %s to directory %s", path, dir)
	http.Handle(path, httpHandler{app, http.StripPrefix(path, http.FileServer(http.Dir(dir)))})
}
