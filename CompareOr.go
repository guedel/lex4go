package main

type CompareOr struct {
	childs []CompareInterface
}

func (c CompareOr) getKind() TestKind {
	return Or
}

func (c CompareOr) accept(g GeneratorInterface) {
	g.VisitCompare(c)
}
