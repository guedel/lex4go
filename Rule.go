package main

import (
	"encoding/xml"

	"github.com/guedel/lex4go/generator"
)

type Rule struct {
	XMLName xml.Name          `xml:"rule"`
	Id      string            `xml:"id,attr"`
	From    string            `xml:"from,attr"`
	To      string            `xml:"to,attr"`
	Test    generator.Compare `xml:"test"`
	Repeat  int               `xml:"repeat,attr"`
	Final   bool              `xml:"final"`
	Concat  bool              `xml:"concat"`
	Reset   bool              `xml:"reset"`
	Action  string            `xml:"action"`
}
