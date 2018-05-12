package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

func writeBuf(caller string, w http.ResponseWriter, buf []byte) {
	_, err := w.Write(buf)
	if err != nil {
		log.Printf("%s writeBuf: %v", caller, err)
	}
}

func writeStr(caller string, w http.ResponseWriter, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		log.Printf("%s writeStr: %v", caller, err)
	}
}

type apiHandler func(w http.ResponseWriter, r *http.Request, app *server)

func registerApi(app *server, path string, handler apiHandler) {
	log.Printf("registering api: %s", path)
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) { handler(w, r, app) })
	app.apis = append(app.apis, path)
}

func serveApiRule(w http.ResponseWriter, r *http.Request, app *server) {
	log.Printf("serveApiRule: url=%s from=%s", r.URL.Path, r.RemoteAddr)

	if authOk := auth(w, r, app); !authOk {
		return
	}

	switch r.Method {
	case http.MethodGet:
		ruleGet(w, r, app)
	case http.MethodPost:
		rulePost(w, r, app)
	default:
		http.Error(w, "Method not supported", 405)
	}
}

func ruleGet(w http.ResponseWriter, r *http.Request, app *server) {

	name := strings.TrimPrefix(r.URL.Path, "/api/rule/")

	if name != "" {
		r, errGet := app.ruleGet(name)
		if errGet != nil {
			log.Printf("ruleGet: not found name=%s: %v", name, errGet)
			http.Error(w, "Not found", 404)
			return
		}

		out, errMarshal := yaml.Marshal(r)
		if errMarshal != nil {
			log.Printf("ruleGet: marshal: %v", errMarshal)
			http.Error(w, "Internal server error", 500)
			return
		}

		writeBuf("ruleGet", w, out)

		return
	}

	out, errDump := app.ruleDump()
	if errDump != nil {
		log.Printf("ruleGet: ruleDump: %v", errDump)
		http.Error(w, "Internal server error", 500)
		return
	}

	writeBuf("ruleGet", w, out)
}

func rulePost(w http.ResponseWriter, r *http.Request, app *server) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("rulePost: body: %v", err)
		http.Error(w, "Internal server error", 500)
		return
	}

	writeBuf("rulePost", w, body)
}

func serveApi(w http.ResponseWriter, r *http.Request, app *server) {
	log.Printf("serveApi: url=%s from=%s", r.URL.Path, r.RemoteAddr)

	if authOk := auth(w, r, app); !authOk {
		return
	}

	apis := app.apiList()
	rules := app.ruleList()

	writeStr("serveApi", w, `<!DOCTYPE html>
<html lang="en-US">
<head>
<title>disbalance api</title>
</head>
<body>
<p>welcome to disbalance api</p>
<a href="/console">console</a>

<p>APIs:</p>
<ul>
`)

	for _, a := range apis {
		writeStr("serveApi", w, fmt.Sprintf(`<li><a href="%s">%s</a></li>`, a, a)+"\n")
	}

	writeStr("serveApi", w, `</ul>
<p>Rules:</p>
<ul>
`)
	for _, r := range rules {
		writeStr("serveApi", w, fmt.Sprintf(`<li><a href="/api/rule/%s">/api/rule/%s</a></li>`, r.Name, r.Name)+"\n")
	}

	writeStr("serveApi", w, `</ul>
</body>
</html>`)

}
