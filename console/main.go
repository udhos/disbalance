package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"

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
	addText.Cols = 21
	addTextDiv1.SetTextContent("rule name")
	addProtoDiv1.SetTextContent("protocol")
	addListenDiv1.SetTextContent("listener")
	addProtoOpt1.Value = "tcp"
	addProtoOpt1.Text = "tcp"
	addProtoOpt2.Value = "udp"
	addProtoOpt2.Text = "udp"
	addListen.Rows = 1
	addListen.MaxLength = 20
	addListen.Cols = 21

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

		col1.SetClass("rule-col1 rule-delete-button")
		col2.SetClass("rule-col2 cell")
		col3.SetClass("rule-col3 cell")
		col4.SetClass("rule-col4 cell")
		col5.SetClass("rule-col5 cell")

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

		ruleListen := func(e dom.Event) {
			url := "/api/rule/" + name
			log.Printf("ruleListen: rule=%s %v url=%s", name, e, url)

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

				t := e.Target().(*dom.HTMLTextAreaElement)
				listenNew := t.Value

				log.Printf("ruleListen: old=%s new=%s", old.Listener, listenNew)
				old.Listener = listenNew

				bufNew, errMarshal := yaml.Marshal(old)
				if errMarshal != nil {
					log.Printf("ruleListen: marshal error: %v", errMarshal)
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
		editListen.AddEventListener("change", false, ruleListen)

		s2.SetTextContent(ruleName)
		s3.AppendChild(editProto)
		s4.AppendChild(editListen)

		// sort targets by name
		targetList := make([]string, 0, len(r.Targets))
		for tn := range r.Targets {
			targetList = append(targetList, tn)
		}
		sort.Strings(targetList)

		targetAddDiv := d.CreateElement("div").(*dom.HTMLDivElement)

		targetTab := d.CreateElement("div").(*dom.HTMLDivElement)
		targetTab.SetClass("table")

		c1 := d.CreateElement("div").(*dom.HTMLDivElement)
		c2 := d.CreateElement("div").(*dom.HTMLDivElement)
		c3 := d.CreateElement("div").(*dom.HTMLDivElement)
		c4 := d.CreateElement("div").(*dom.HTMLDivElement)
		c5 := d.CreateElement("div").(*dom.HTMLDivElement)
		c6 := d.CreateElement("div").(*dom.HTMLDivElement)

		c1.SetClass("line")
		c2.SetClass("line")
		c3.SetClass("line")
		c4.SetClass("line")
		c5.SetClass("line")
		c6.SetClass("line")

		h2 := d.CreateElement("div").(*dom.HTMLDivElement)
		h3 := d.CreateElement("div").(*dom.HTMLDivElement)
		h4 := d.CreateElement("div").(*dom.HTMLDivElement)
		h5 := d.CreateElement("div").(*dom.HTMLDivElement)
		h6 := d.CreateElement("div").(*dom.HTMLDivElement)

		h2.SetTextContent("target")
		h3.SetTextContent("interval")
		h4.SetTextContent("timeout")
		h5.SetTextContent("minimum")
		h6.SetTextContent("address")

		maxLen := 30
		cols := 4

		text2 := d.CreateElement("textarea").(*dom.HTMLTextAreaElement)
		text2.Rows = 1
		text2.MaxLength = maxLen
		text2.Cols = cols
		text3 := d.CreateElement("textarea").(*dom.HTMLTextAreaElement)
		text3.Rows = 1
		text3.MaxLength = maxLen
		text3.Cols = cols
		text4 := d.CreateElement("textarea").(*dom.HTMLTextAreaElement)
		text4.Rows = 1
		text4.MaxLength = maxLen
		text4.Cols = cols
		text5 := d.CreateElement("textarea").(*dom.HTMLTextAreaElement)
		text5.Rows = 1
		text5.MaxLength = maxLen
		text5.Cols = cols
		text6 := d.CreateElement("textarea").(*dom.HTMLTextAreaElement)
		text6.Rows = 1
		text6.MaxLength = maxLen
		text6.Cols = cols

		textDiv2 := d.CreateElement("div").(*dom.HTMLDivElement)
		textDiv3 := d.CreateElement("div").(*dom.HTMLDivElement)
		textDiv4 := d.CreateElement("div").(*dom.HTMLDivElement)
		textDiv5 := d.CreateElement("div").(*dom.HTMLDivElement)
		textDiv6 := d.CreateElement("div").(*dom.HTMLDivElement)

		textDiv2.AppendChild(text2)
		textDiv3.AppendChild(text3)
		textDiv4.AppendChild(text4)
		textDiv5.AppendChild(text5)
		textDiv6.AppendChild(text6)

		targetAddImg := d.CreateElement("img").(*dom.HTMLImageElement)
		targetAddImg.Src = "/console/plus.png"
		targetAddImg.Height = 16
		targetAddImg.Width = 16
		targetAddBut := d.CreateElement("button").(*dom.HTMLButtonElement)
		targetAddBut.SetClass("unstyled-button")
		targetAddBut.AppendChild(targetAddImg)

		c1.AppendChild(targetAddBut)
		c2.AppendChild(h2)
		c2.AppendChild(textDiv2)
		c3.AppendChild(h3)
		c3.AppendChild(textDiv3)
		c4.AppendChild(h4)
		c4.AppendChild(textDiv4)
		c5.AppendChild(h5)
		c5.AppendChild(textDiv5)
		c6.AppendChild(h6)
		c6.AppendChild(textDiv6)

		targetAddDiv.AppendChild(c1)
		targetAddDiv.AppendChild(c2)
		targetAddDiv.AppendChild(c3)
		targetAddDiv.AppendChild(c4)
		targetAddDiv.AppendChild(c5)
		targetAddDiv.AppendChild(c6)

		s5.AppendChild(targetAddDiv)
		s5.AppendChild(targetTab)

		for _, targetName := range targetList {
			t := r.Targets[targetName]
			targetLine := d.CreateElement("div").(*dom.HTMLDivElement)
			targetLine.SetClass("table-line")

			col1 := d.CreateElement("div").(*dom.HTMLDivElement)
			col2 := d.CreateElement("div").(*dom.HTMLDivElement)

			col1.SetClass("line")
			col2.SetClass("line")

			targetDelBut := d.CreateElement("button").(*dom.HTMLButtonElement)
			targetDelBut.SetClass("unstyled-button")
			targetDelImg := d.CreateElement("img").(*dom.HTMLImageElement)
			targetDelImg.Src = "/console/trash.png"
			targetDelImg.Height = 16
			targetDelImg.Width = 16
			targetDelBut.AppendChild(targetDelImg)

			col1.AppendChild(targetDelBut)
			col2.SetTextContent(fmt.Sprintf("%s: %v", targetName, t))

			targetLine.AppendChild(col1)
			targetLine.AppendChild(col2)

			targetTab.AppendChild(targetLine)
		}

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

		targetAdd := func(e dom.Event) {

			target := text2.Value
			log.Printf("targetAdd: rule=%s target=%s", name, target)

			if target == "" {
				log.Printf("targetAdd: empty target name")
				return
			}

			interval := text3.Value
			timeout := text4.Value
			minimum := text5.Value
			address := text6.Value

			vInt, _ := strconv.Atoi(interval)
			vTmout, _ := strconv.Atoi(timeout)
			vMin, _ := strconv.Atoi(minimum)

			r := rule.Rule{
				Targets: map[string]rule.Target{},
			}
			r.Targets[target] = rule.Target{
				Check: rule.HealthCheck{
					Interval: vInt,
					Timeout:  vTmout,
					Minimum:  vMin,
					Address:  address,
				},
			}

			rules := map[string]rule.Rule{}
			rules[name] = r

			buf, errMarshal := yaml.Marshal(rules)
			if errMarshal != nil {
				log.Printf("targetAdd: marshal: %v", errMarshal)
				return
			}

			// goroutine needed to prevent block
			go func() {
				_, errPost := httpPost("/api/rule/", "application/x-yaml", buf)
				if errPost != nil {
					log.Printf("targetAdd: rule=%s target=%s: error: %v", name, target, errPost)
					return
				}
				log.Printf("targetAdd: rule=%s target=%s added", name, target)

				loadRules()
			}()
		}

		targetAddBut.AddEventListener("click", false, targetAdd)

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
