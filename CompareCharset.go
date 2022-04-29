package main

import (
	"encoding/xml"
)

type CompareCharset struct {
	XMLName xml.Name `xml:"charset"`
	Name    string   `xml:",chardata"`
}

func (c CompareCharset) getKind() TestKind {
	return Charset
}

func (c CompareCharset) accept(g GeneratorInterface) {
	g.VisitCompare(c)
}

/*
func (c *CompareCharset) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	fmt.Println("CompareCharset.UnmarshalXML: ", start)
	return nil
}
*/
