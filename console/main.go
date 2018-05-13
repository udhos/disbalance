package main

import (
	"fmt"

	//"github.com/gopherjs/gopherjs"
	"gopkg.in/yaml.v2"
	"honnef.co/go/js/dom"

	"github.com/udhos/disbalance/rule"
)

func main() {
	println("main: hello from console")

	loadRules()
}

func loadRules() {
	d := dom.GetWindow().Document()
	div := d.GetElementByID("rules").(*dom.HTMLDivElement)
	div.SetTextContent("loading rules")

	buf, errFetch := httpFetch("/api/rule/")
	if errFetch != nil {
		div.SetTextContent("fetch fail: " + errFetch.Error())
		return
	}

	var rules []rule.Rule
	if errYaml := yaml.Unmarshal(buf, &rules); errYaml != nil {
		div.SetTextContent("yaml fail: " + errYaml.Error())
		return
	}

	removeChildren(div.BasicNode)

	for _, r := range rules {
		line := d.CreateElement("div").(*dom.HTMLDivElement)
		but := d.CreateElement("button").(*dom.HTMLButtonElement)
		img := d.CreateElement("img").(*dom.HTMLImageElement)
		span1 := d.CreateElement("span").(*dom.HTMLSpanElement)
		span2 := d.CreateElement("span").(*dom.HTMLSpanElement)

		but.SetClass("unstyled-button")
		img.Src = "/console/trash.png"
		img.Height = 16
		img.Width = 16
		span2.SetTextContent(fmt.Sprintf("%v", r))

		span1.AppendChild(img)
		but.AppendChild(span1)
		line.AppendChild(but)
		line.AppendChild(span2)

		div.AppendChild(line)
	}
}

func removeChildren(n *dom.BasicNode) {
	for _, c := range n.ChildNodes() {
		n.RemoveChild(c)
	}
}
