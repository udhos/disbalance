package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/udhos/disbalance/rule"
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

// https://stackoverflow.com/questions/630453/put-vs-post-in-rest
//
// POST to a URL creates a child resource at a server defined URL.
// PUT to a URL creates/replaces the resource in its entirety at the client defined URL.
// PATCH to a URL updates part of the resource at that client defined URL.

func serveApiRule(w http.ResponseWriter, r *http.Request, app *server) {
	log.Printf("serveApiRule: url=%s from=%s", r.URL.Path, r.RemoteAddr)

	if authOk := auth(w, r, app); !authOk {
		return
	}

	switch r.Method {
	case http.MethodGet:
		ruleGet(w, r, app)
	case http.MethodDelete:
		ruleDelete(w, r, app)
	case http.MethodPost:
		rulePost(w, r, app)
	case http.MethodPut:
		rulePut(w, r, app)
	default:
		http.Error(w, "Method not supported", 405)
	}
}

func ruleDelete(w http.ResponseWriter, r *http.Request, app *server) {

	name := strings.TrimPrefix(r.URL.Path, "/api/rule/")
	if name == "" {
		log.Printf("ruleDelete: missing name")
		http.Error(w, "Missing rule name", 400)
		return
	}

	app.ruleDel(name)

	http.Error(w, "Ok", 200)
}

func rulePut(w http.ResponseWriter, r *http.Request, app *server) {

	name := strings.TrimPrefix(r.URL.Path, "/api/rule/")
	if name == "" {
		log.Printf("rulePut: missing name")
		http.Error(w, "Missing rule name", 400)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("rulePut: body: %v", err)
		http.Error(w, "Internal server error", 500)
		return
	}

	var ruleSingle rule.Rule

	if errYaml := yaml.Unmarshal(body, &ruleSingle); errYaml != nil {
		log.Printf("rulePut: unmarshal: %v", errYaml)
		http.Error(w, "Bad request", 400)
		return
	}

	ruleSingle.ForceValidChecks()

	app.rulePut(name, ruleSingle)

	http.Error(w, fmt.Sprintf("Rule updated: %s", name), 200)
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

	rules := map[string]rule.Rule{}

	if errYaml := yaml.Unmarshal(body, &rules); errYaml != nil {
		log.Printf("rulePost: unmarshal: %v", errYaml)
		http.Error(w, "Bad request", 400)
		return
	}

	for _, r := range rules {
		r.ForceValidChecks()
	}

	app.rulePost(rules)

	http.Error(w, fmt.Sprintf("Rules updated: %d", len(rules)), 200)
}

func serveApi(w http.ResponseWriter, r *http.Request, app *server) {
	log.Printf("serveApi: url=%s from=%s", r.URL.Path, r.RemoteAddr)

	if authOk := auth(w, r, app); !authOk {
		return
	}

	apis := app.apiList()
	rules := app.ruleTable()

	writeStr("serveApi", w, `<!DOCTYPE html>
<html lang="en-US">
<head>
<title>disbalance api</title>
</head>
<body>
<p>welcome to disbalance api - <a href="https://github.com/udhos/disbalance">github</a></p>
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
<table border>
<thead>
<th>Rule</th>
<th>Protocol</th>
<th>Listener</th>
<th>Targets</th>
</thead>
<tbody>
`)
	for name := range rules {
		r := rules[name]

		writeStr("serveApi", w, fmt.Sprintf(`
<tr>
<td><a href="/api/rule/%s">/api/rule/%s</a></td>
`, name, name))

		writeStr("serveApi", w, "<td>")
		writeStr("serveApi", w, r.Protocol)
		writeStr("serveApi", w, "</td>\n")

		writeStr("serveApi", w, "<td>")
		writeStr("serveApi", w, r.Listener)
		writeStr("serveApi", w, "</td>\n")

		writeStr("serveApi", w, "<td>")
		for a, t := range r.Targets {
			writeStr("serveApi", w, "<div>")
			writeStr("serveApi", w, fmt.Sprintf("%s interval=%d timeout=%d minimum=%d address=%s", a, t.Check.Interval, t.Check.Timeout, t.Check.Minimum, t.Check.Address))
			writeStr("serveApi", w, "</div>")
		}
		writeStr("serveApi", w, "</td>\n")

		writeStr("serveApi", w, "</tr>\n")
	}

	writeStr("serveApi", w, `</ul>
</tbody>
</table>
</body>
</html>`)

}
