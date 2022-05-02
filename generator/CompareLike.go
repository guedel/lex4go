package generator

import "encoding/xml"

type CompareLike struct {
	XMLName    xml.Name `xml:"like"`
	Expression string   `xml:",chardata"`
}

func (c CompareLike) getKind() TestKind {
	return Like
}

func (c CompareLike) accept(g GeneratorInterface) {
	g.VisitCompare(c)
}
