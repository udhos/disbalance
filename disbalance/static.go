package main

import (
	"log"
	"net/http"
)

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
