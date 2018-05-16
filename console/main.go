package main

import (
	"fmt"

	//"github.com/gopherjs/gopherjs"
	"gopkg.in/yaml.v2"
	"honnef.co/go/js/dom"

	"github.com/udhos/disbalance/rule"
)

func main() {
	println("hello from console")

	buf, errFetch := httpFetch("/api/rule/")
	if errFetch != nil {
		println("fail: " + errFetch.Error())
		return
	}

	var rules []rule.Rule

	if errYaml := yaml.Unmarshal(buf, &rules); errYaml != nil {
		println("fail: " + errYaml.Error())
		return
	}

	d := dom.GetWindow().Document()
	div := d.GetElementByID("rules").(*dom.HTMLDivElement)
	for _, r := range rules {
		s := fmt.Sprintf("%v", r)
		child := d.CreateElement("div").(*dom.HTMLDivElement)
		child.SetTextContent(s)
		div.AppendChild(child)
	}
}
