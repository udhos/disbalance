package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"golang.org/x/crypto/acme/autocert"
)

const (
	version = "0.0"
)

type config struct {
}

type server struct {
	cfg config

	config        string
	controlPort   string
	controlDomain string
}

func main() {

	var app server

	log.Printf("version %s runtime %s", version, runtime.Version())

	flag.StringVar(&app.config, "config", "/etc/disbalance.conf", "config path")
	flag.StringVar(&app.controlPort, "controlPort", ":8080", "control port")
	flag.StringVar(&app.controlDomain, "controlDomain", "example.com", "control certificate domain")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, TLS user! Your config: %+v", r.TLS)
	})
	log.Fatal(http.Serve(autocert.NewListener(app.controlDomain+app.controlPort), mux))
}
