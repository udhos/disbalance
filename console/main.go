package main

import (
	"fmt"
	"log"
	"sort"

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

	rules := map[string]rule.Rule{}
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
	addSpan.SetClass("line")
	addTextSpan.SetClass("line")
	addProtoSpan.SetClass("line")
	addListenSpan.SetClass("line")

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

	// sort rules by name
	ruleList := make([]string, 0, len(rules))
	for n := range rules {
		ruleList = append(ruleList, n)
	}
	sort.Strings(ruleList)

	tab := d.CreateElement("div").(*dom.HTMLDivElement)
	tab.SetClass("table")
	div.AppendChild(tab)

	for _, ruleName := range ruleList {
		r := rules[ruleName]

		line := d.CreateElement("div").(*dom.HTMLDivElement)
		line.SetClass("table-line")

		but := d.CreateElement("button").(*dom.HTMLButtonElement)
		img := d.CreateElement("img").(*dom.HTMLImageElement)

		col1 := d.CreateElement("div").(*dom.HTMLDivElement)
		col2 := d.CreateElement("div").(*dom.HTMLDivElement)
		col3 := d.CreateElement("div").(*dom.HTMLDivElement)
		col4 := d.CreateElement("div").(*dom.HTMLDivElement)
		col5 := d.CreateElement("div").(*dom.HTMLDivElement)

		col1.SetClass("rule-delete-button")
		col2.SetClass("cell")
		col3.SetClass("cell")
		col4.SetClass("cell")
		col5.SetClass("cell")

		s2 := d.CreateElement("span").(*dom.HTMLSpanElement)
		s3 := d.CreateElement("span").(*dom.HTMLSpanElement)
		s4 := d.CreateElement("span").(*dom.HTMLSpanElement)
		s5 := d.CreateElement("span").(*dom.HTMLSpanElement)

		name := ruleName // save name for closure below

		ruleDelete := func(e dom.Event) {
			log.Printf("ruleDelete: rule=%s %v", name, e)

			// goroutine needed to prevent block
			go func() {
				_, errDel := httpDelete("/api/rule/" + name)
				if errDel != nil {
					log.Printf("delete error: %v", errDel)
					return
				}
				log.Printf("deleted: %s", name)

				loadRules()
			}()
		}

		ruleProto := func(e dom.Event) {
			url := "/api/rule/" + name
			log.Printf("ruleProto: rule=%s %v url=%s", name, e, url)

			// goroutine needed to prevent block
			go func() {
				// fetch old rule
				buf, errFetch := httpFetch(url)
				if errFetch != nil {
					log.Printf("fetch fail: " + errFetch.Error())
					return
				}

				var old rule.Rule
				if errYaml := yaml.Unmarshal(buf, &old); errYaml != nil {
					log.Printf("yaml fail: " + errYaml.Error())
					return
				}

				t := e.Target()
				s := t.(*dom.HTMLSelectElement)
				protoNew := s.Value

				log.Printf("ruleProto: old=%s new=%s", old.Protocol, protoNew)
				old.Protocol = protoNew

				bufNew, errMarshal := yaml.Marshal(old)
				if errMarshal != nil {
					log.Printf("ruleProto: marshal error: %v", errMarshal)
					return
				}

				_, errPut := httpPut(url, "application/x-yaml", bufNew)
				if errPut != nil {
					log.Printf("put error: %v", errPut)
					return
				}

				//loadRules()
			}()
		}

		but.SetClass("unstyled-button")
		but.AddEventListener("click", false, ruleDelete)
		img.Src = "/console/trash.png"
		img.Height = 16
		img.Width = 16

		editProto := d.CreateElement("select").(*dom.HTMLSelectElement)
		editProtoOpt1 := d.CreateElement("option").(*dom.HTMLOptionElement)
		editProtoOpt1.Value = "tcp"
		editProtoOpt1.Text = "tcp"
		editProtoOpt2 := d.CreateElement("option").(*dom.HTMLOptionElement)
		editProtoOpt2.Value = "udp"
		editProtoOpt2.Text = "udp"
		editProto.AppendChild(editProtoOpt1)
		editProto.AppendChild(editProtoOpt2)
		if r.Protocol == "tcp" {
			editProto.SelectedIndex = 0
		} else {
			editProto.SelectedIndex = 1
		}
		editProto.AddEventListener("change", false, ruleProto)
		editListen := d.CreateElement("textarea").(*dom.HTMLTextAreaElement)
		editListen.Rows = 1
		editListen.Value = r.Listener

		s2.SetTextContent(ruleName)
		s3.AppendChild(editProto)
		s4.AppendChild(editListen)
		s5.SetTextContent(fmt.Sprintf("%v", r.Targets))

		col2.AppendChild(s2)
		col3.AppendChild(s3)
		col4.AppendChild(s4)
		col5.AppendChild(s5)

		but.AppendChild(img)
		col1.AppendChild(but)

		line.AppendChild(col1)
		line.AppendChild(col2)
		line.AppendChild(col3)
		line.AppendChild(col4)
		line.AppendChild(col5)

		tab.AppendChild(line)
	}

	ruleAdd := func(e dom.Event) {

		name := addText.Value
		if name == "" {
			log.Printf("add: empty name")
			return
		}

		newRule := rule.Rule{
			Protocol: addProto.Value,
			Listener: addListen.Value,
		}
		// addText addProto addListen

		ruleTab := map[string]rule.Rule{
			name: newRule,
		}

		body, errMarshal := yaml.Marshal(ruleTab)
		if errMarshal != nil {
			log.Printf("add marshal error: %v", errMarshal)
			return
		}

		// goroutine needed to prevent block
		go func() {
			_, errAdd := httpPost("/api/rule/", "application/x-yaml", body)
			if errAdd != nil {
				log.Printf("add error: %v", errAdd)
				return
			}
			log.Printf("added: %s", name)

			loadRules()
		}()
	}

	addBut.AddEventListener("click", false, ruleAdd)
}

func removeChildren(n *dom.BasicNode) {
	for _, c := range n.ChildNodes() {
		n.RemoveChild(c)
	}
}
