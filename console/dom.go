package main

import (
	"honnef.co/go/js/dom"
)

func newInputText(d dom.Document, parent *dom.BasicNode, size, maxLength int) *dom.HTMLInputElement {

	input := d.CreateElement("input").(*dom.HTMLInputElement)

	input.Size = size
	input.MaxLength = maxLength
	input.Type = "text"

	parent.AppendChild(input)

	return input
}
