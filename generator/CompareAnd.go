package generator

type CompareAnd struct {
	childs []CompareInterface
}

func (c CompareAnd) getKind() TestKind {
	return And
}

func (c CompareAnd) accept(g GeneratorInterface) {
	g.VisitCompare(c)
}
