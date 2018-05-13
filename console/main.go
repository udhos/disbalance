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

	addDiv := d.CreateElement("div").(*dom.HTMLDivElement)
	addSpan := d.CreateElement("span").(*dom.HTMLSpanElement)
	addBut := d.CreateElement("button").(*dom.HTMLButtonElement)
	addImg := d.CreateElement("img").(*dom.HTMLImageElement)
	addText := d.CreateElement("textarea").(*dom.HTMLTextAreaElement)
	addTextSpan := d.CreateElement("span").(*dom.HTMLSpanElement)
	addTextDiv1 := d.CreateElement("div").(*dom.HTMLDivElement)
	addTextDiv2 := d.CreateElement("div").(*dom.HTMLDivElement)

	addBut.SetClass("unstyled-button")
	addSpan.SetClass("inline")
	addTextSpan.SetClass("inline")
	addImg.Src = "/console/plus.png"
	addImg.Height = 16
	addImg.Width = 16
	addText.Rows = 1
	addText.MaxLength = 20
	addTextDiv1.SetTextContent("rule name")

	addSpan.AppendChild(addBut)
	addTextDiv2.AppendChild(addText)
	addTextSpan.AppendChild(addTextDiv1)
	addTextSpan.AppendChild(addTextDiv2)
	addBut.AppendChild(addImg)
	addDiv.AppendChild(addSpan)
	addDiv.AppendChild(addTextSpan)
	div.AppendChild(addDiv)

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
