package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

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

	out, errList := app.ruleList()
	if errList != nil {
		log.Printf("serveApiRuleList: ruleList: %v", errList)
	}

	_, errWrite := w.Write(out)
	if errWrite != nil {
		log.Printf("serveApiRuleList: write: %v", errWrite)
	}
}

func serveApi(w http.ResponseWriter, r *http.Request, app *server) {
	log.Printf("serveApi: url=%s from=%s", r.URL.Path, r.RemoteAddr)

	if authOk := auth(w, r, app); !authOk {
		return
	}

	apis := app.apiList()

	io.WriteString(w, `<!DOCTYPE html>
<html lang="en-US">
<title>
<head>disbalance api</head>
</title>
<body>
<p>welcome to the api</p>
<a href="/console">console</a>
<ul>`)

	for _, a := range apis {
		io.WriteString(w, fmt.Sprintf(`<li><a href="%s">%s</a></li>`, a, a))
	}

	io.WriteString(w, `</ul>
</body>
</html>`)

}
