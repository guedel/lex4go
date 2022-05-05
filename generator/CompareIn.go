package generator

import "encoding/xml"

type CompareIn struct {
	XMLName    xml.Name `xml:"in"`
	Expression string   `xml:",chardata"`
}

func (c CompareIn) getKind() TestKind {
	return In
}

func (c CompareIn) accept(g GeneratorInterface) {
	g.VisitCompare(c)
}
