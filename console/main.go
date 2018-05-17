package main

import (
	"fmt"
	"log"

	//"github.com/gopherjs/gopherjs"
	"gopkg.in/yaml.v2"
	"honnef.co/go/js/dom"

	"github.com/udhos/disbalance/rule"
)

func main() {
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
	addProto := d.CreateElement("select").(*dom.HTMLSelectElement)
	addProtoSpan := d.CreateElement("span").(*dom.HTMLSpanElement)
	addProtoDiv1 := d.CreateElement("div").(*dom.HTMLDivElement)
	addProtoDiv2 := d.CreateElement("div").(*dom.HTMLDivElement)
	addProtoOpt1 := d.CreateElement("option").(*dom.HTMLOptionElement)
	addProtoOpt2 := d.CreateElement("option").(*dom.HTMLOptionElement)
	addListen := d.CreateElement("textarea").(*dom.HTMLTextAreaElement)
	addListenSpan := d.CreateElement("span").(*dom.HTMLSpanElement)
	addListenDiv1 := d.CreateElement("div").(*dom.HTMLDivElement)
	addListenDiv2 := d.CreateElement("div").(*dom.HTMLDivElement)

	addBut.SetClass("unstyled-button")
	addSpan.SetClass("inline")
	addTextSpan.SetClass("inline")
	addProtoSpan.SetClass("inline")
	addListenSpan.SetClass("inline")
	addImg.Src = "/console/plus.png"
	addImg.Height = 16
	addImg.Width = 16
	addText.Rows = 1
	addText.MaxLength = 20
	addTextDiv1.SetTextContent("rule name")
	addProtoDiv1.SetTextContent("protocol")
	addListenDiv1.SetTextContent("listener")
	addProtoOpt1.Value = "tcp"
	addProtoOpt1.Text = "tcp"
	addProtoOpt2.Value = "udp"
	addProtoOpt2.Text = "udp"
	addListen.Rows = 1
	addListen.MaxLength = 20

	addSpan.AppendChild(addBut)
	addTextDiv2.AppendChild(addText)
	addTextSpan.AppendChild(addTextDiv1)
	addTextSpan.AppendChild(addTextDiv2)
	addProto.AppendChild(addProtoOpt1)
	addProto.AppendChild(addProtoOpt2)
	addProtoDiv2.AppendChild(addProto)
	addProtoSpan.AppendChild(addProtoDiv1)
	addProtoSpan.AppendChild(addProtoDiv2)
	addListenDiv2.AppendChild(addListen)
	addListenSpan.AppendChild(addListenDiv1)
	addListenSpan.AppendChild(addListenDiv2)
	addBut.AppendChild(addImg)
	addDiv.AppendChild(addSpan)
	addDiv.AppendChild(addTextSpan)
	addDiv.AppendChild(addProtoSpan)
	addDiv.AppendChild(addListenSpan)
	div.AppendChild(addDiv)

	for _, r := range rules {
		line := d.CreateElement("div").(*dom.HTMLDivElement)
		but := d.CreateElement("button").(*dom.HTMLButtonElement)
		img := d.CreateElement("img").(*dom.HTMLImageElement)
		span1 := d.CreateElement("span").(*dom.HTMLSpanElement)
		span2 := d.CreateElement("span").(*dom.HTMLSpanElement)

		ruleName := r.Name

		ruleDelete := func(e dom.Event) {
			log.Printf("ruleDelete: rule=%s %v", ruleName, e)
		}

		but.SetClass("unstyled-button")
		but.AddEventListener("click", false, ruleDelete)
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
