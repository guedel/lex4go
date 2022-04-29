package main

type CompareCharset struct {
	name string
}

func (c CompareCharset) getKind() TestKind {
	return Charset
}

func (c CompareCharset) accept(g GeneratorInterface) {
	g.VisitCompare(c)
}
